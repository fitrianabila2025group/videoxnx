package routes

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/config"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/controllers"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/middleware"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fmw "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

// NewApp wires routes and middleware.
func NewApp(db *gorm.DB, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "videoxnx-backend",
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
	})

	app.Use(recover.New())
	app.Use(fmw.New())
	app.Use(middleware.SecureHeaders())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSAllowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: contains(cfg.CORSAllowedOrigins, "*") == false,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        120,
		Expiration: 1 * time.Minute,
	}))

	pub := &controllers.Public{DB: db}
	adm := &controllers.Admin{DB: db, Cfg: cfg}

	app.Get("/healthz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"ok": true}) })

	api := app.Group("/api")
	// Public
	api.Get("/posts", pub.ListPosts)
	api.Get("/posts/:slug", pub.GetPost)
	api.Get("/categories", pub.ListCategories)
	api.Get("/categories/:slug/posts", pub.PostsByCategory)
	api.Get("/tags", pub.ListTags)
	api.Get("/tags/:slug/posts", pub.PostsByTag)
	api.Get("/search", pub.Search)
	api.Get("/trending", pub.Trending)
	api.Post("/reports", pub.CreateReport)
	api.Get("/img", pub.ImageProxy)
	api.Get("/site", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"site_name":        cfg.SiteName,
			"meta_title":       cfg.MetaTitle,
			"meta_description": cfg.MetaDescription,
			"age_gate_enabled": cfg.AgeGateEnabled,
			"dmca_email":       cfg.DMCAEmail,
		})
	})

	// Admin
	admin := api.Group("/admin")
	admin.Post("/login", adm.Login)
	admin.Post("/logout", adm.Logout)

	priv := admin.Use(middleware.RequireAdmin(cfg))
	priv.Get("/dashboard", adm.Dashboard)

	priv.Get("/posts", adm.ListPosts)
	priv.Put("/posts/:id", adm.UpdatePost)
	priv.Delete("/posts/:id", adm.DeletePost)
	priv.Patch("/posts/:id/status", adm.UpdatePostStatus)

	priv.Get("/categories", adm.ListCategories)
	priv.Post("/categories", adm.CreateCategory)
	priv.Put("/categories/:id", adm.UpdateCategory)
	priv.Delete("/categories/:id", adm.DeleteCategory)

	priv.Get("/tags", adm.ListTags)
	priv.Post("/tags", adm.CreateTag)
	priv.Put("/tags/:id", adm.UpdateTag)
	priv.Delete("/tags/:id", adm.DeleteTag)

	priv.Post("/scraper/run", adm.ScraperRun)
	priv.Get("/scraper/logs", adm.ScraperLogs)

	priv.Get("/reports", adm.ListReports)
	priv.Patch("/reports/:id", adm.UpdateReport)

	priv.Get("/settings", adm.ListSettings)
	priv.Put("/settings", adm.UpdateSettings)

	// SEO
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		c.Type("txt")
		return c.SendString("User-agent: *\nAllow: /\nSitemap: " + cfg.AppURL + "/sitemap.xml\n")
	})
	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		return sitemap(c, db, cfg)
	})

	// Frontend proxy: when FRONTEND_URL is set, forward all unmatched
	// requests to the Next.js process (used by the all-in-one Docker image
	// where Next.js runs on 127.0.0.1:3000 alongside this server).
	if cfg.FrontendURL != "" {
		target := strings.TrimRight(cfg.FrontendURL, "/")
		app.Use(func(c *fiber.Ctx) error {
			url := target + c.OriginalURL()
			if err := proxy.Do(c, url); err != nil {
				return fiber.NewError(fiber.StatusBadGateway, "frontend unavailable")
			}
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		})
	}

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := err.Error()
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

type sitemapURL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod,omitempty"`
}
type sitemapSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

func sitemap(c *fiber.Ctx, db *gorm.DB, cfg *config.Config) error {
	var posts []models.Post
	db.Where("status = ? AND safety_status <> ?", "published", "blocked").
		Order("COALESCE(published_at, scraped_at) DESC").Limit(2000).Find(&posts)
	urls := []sitemapURL{
		{Loc: cfg.AppURL + "/"},
		{Loc: cfg.AppURL + "/latest"},
		{Loc: cfg.AppURL + "/trending"},
		{Loc: cfg.AppURL + "/dmca"},
	}
	for _, p := range posts {
		u := sitemapURL{Loc: fmt.Sprintf("%s/post/%s", cfg.AppURL, p.Slug)}
		if p.PublishedAt != nil {
			u.LastMod = p.PublishedAt.Format("2006-01-02")
		}
		urls = append(urls, u)
	}
	out, err := xml.MarshalIndent(sitemapSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: urls}, "", "  ")
	if err != nil {
		return err
	}
	c.Type("xml")
	return c.Send(append([]byte(xml.Header), out...))
}

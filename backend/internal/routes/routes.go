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
		// Trust the platform proxy (Railway, Render, Fly, Nginx, ...) so that
		// Ctx.IP() returns the real client IP. Without this every visitor
		// appears as the same edge IP and the rate limiter throttles the
		// whole site after ~120 requests/min.
		ProxyHeader:             fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: false,
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

	pub := &controllers.Public{DB: db}
	adm := &controllers.Admin{DB: db, Cfg: cfg}

	app.Get("/healthz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"ok": true}) })

	api := app.Group("/api")

	// Per-IP rate limit, scoped to the public API only. The Next.js frontend
	// (served via the proxy below) and admin endpoints are NOT rate-limited
	// here so that legitimate page loads (which trigger many /_next/* sub
	// requests) and the admin login form keep working.
	api.Use(limiter.New(limiter.Config{
		Max:        300,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			// Skip rate limiting for admin endpoints (login, dashboard, ...).
			return strings.HasPrefix(c.Path(), "/api/admin")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(fiber.Map{"error": "rate limit exceeded, please slow down"})
		},
	}))

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

	// SEO. URLs are derived from the incoming request so the same binary
	// works behind any domain (Railway, Render, custom, ...). cfg.SiteURL
	// (env SITE_URL) is used as a fallback for cases where the request has
	// no Host header (extremely rare).
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		base := publicBaseURL(c, cfg)
		c.Type("txt")
		return c.SendString("User-agent: *\nAllow: /\nSitemap: " + base + "/sitemap.xml\n")
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
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod,omitempty"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}
type sitemapSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// publicBaseURL returns the canonical, public base URL (scheme + host) for the
// current request. It honours X-Forwarded-* headers set by the platform proxy
// (Railway, Render, Fly, Cloudflare, ...). Falls back to cfg.SiteURL, then
// cfg.AppURL. The result never has a trailing slash.
func publicBaseURL(c *fiber.Ctx, cfg *config.Config) string {
	host := strings.TrimSpace(c.Get("X-Forwarded-Host"))
	if host == "" {
		host = strings.TrimSpace(c.Hostname())
	}
	if host != "" && !isLocalHost(host) {
		proto := strings.TrimSpace(c.Get("X-Forwarded-Proto"))
		if proto == "" {
			proto = c.Protocol()
		}
		if proto == "" {
			proto = "https"
		}
		return strings.TrimRight(proto+"://"+host, "/")
	}
	if cfg.SiteURL != "" {
		return strings.TrimRight(cfg.SiteURL, "/")
	}
	if host != "" {
		proto := c.Protocol()
		if proto == "" {
			proto = "http"
		}
		return strings.TrimRight(proto+"://"+host, "/")
	}
	return strings.TrimRight(cfg.AppURL, "/")
}

func isLocalHost(host string) bool {
	h := strings.ToLower(host)
	if i := strings.IndexByte(h, ':'); i >= 0 {
		h = h[:i]
	}
	return h == "localhost" || h == "127.0.0.1" || h == "0.0.0.0" || h == "::1"
}

func sitemap(c *fiber.Ctx, db *gorm.DB, cfg *config.Config) error {
	base := publicBaseURL(c, cfg)

	var posts []models.Post
	db.Where("status = ? AND safety_status <> ?", "published", "blocked").
		Order("COALESCE(published_at, scraped_at) DESC").Limit(5000).Find(&posts)

	urls := []sitemapURL{
		{Loc: base + "/", ChangeFreq: "hourly", Priority: "1.0"},
		{Loc: base + "/latest", ChangeFreq: "hourly", Priority: "0.9"},
		{Loc: base + "/trending", ChangeFreq: "daily", Priority: "0.8"},
		{Loc: base + "/categories", ChangeFreq: "weekly", Priority: "0.6"},
		{Loc: base + "/tags", ChangeFreq: "weekly", Priority: "0.5"},
		{Loc: base + "/dmca", ChangeFreq: "yearly", Priority: "0.2"},
		{Loc: base + "/contact", ChangeFreq: "yearly", Priority: "0.2"},
		{Loc: base + "/disclaimer", ChangeFreq: "yearly", Priority: "0.2"},
		{Loc: base + "/privacy", ChangeFreq: "yearly", Priority: "0.2"},
		{Loc: base + "/age-verification", ChangeFreq: "yearly", Priority: "0.2"},
	}

	// Categories
	var cats []models.Category
	db.Order("name ASC").Limit(2000).Find(&cats)
	for _, ct := range cats {
		if ct.Slug == "" {
			continue
		}
		urls = append(urls, sitemapURL{
			Loc:        fmt.Sprintf("%s/category/%s", base, ct.Slug),
			ChangeFreq: "daily",
			Priority:   "0.6",
		})
	}

	// Tags
	var tags []models.Tag
	db.Order("name ASC").Limit(5000).Find(&tags)
	for _, tg := range tags {
		if tg.Slug == "" {
			continue
		}
		urls = append(urls, sitemapURL{
			Loc:        fmt.Sprintf("%s/tag/%s", base, tg.Slug),
			ChangeFreq: "weekly",
			Priority:   "0.4",
		})
	}

	// Posts
	for _, p := range posts {
		if p.Slug == "" {
			continue
		}
		u := sitemapURL{
			Loc:        fmt.Sprintf("%s/post/%s", base, p.Slug),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		}
		switch {
		case p.PublishedAt != nil:
			u.LastMod = p.PublishedAt.UTC().Format("2006-01-02")
		case !p.ScrapedAt.IsZero():
			u.LastMod = p.ScrapedAt.UTC().Format("2006-01-02")
		}
		urls = append(urls, u)
	}

	out, err := xml.MarshalIndent(sitemapSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}, "", "  ")
	if err != nil {
		return err
	}
	c.Type("xml")
	return c.Send(append([]byte(xml.Header), out...))
}

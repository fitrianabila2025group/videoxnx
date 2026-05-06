package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/config"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/scraper"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/services"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Admin struct {
	DB        *gorm.DB
	Cfg       *config.Config
	Scheduler *scraper.Scheduler // optional, may be nil
}

// POST /api/admin/login
func (a *Admin) Login(c *fiber.Ctx) error {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	var u models.User
	if err := a.DB.Where("email = ?", in.Email).First(&u).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	if !services.VerifyPassword(u.PasswordHash, in.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	token, err := utils.GenerateJWT(a.Cfg.JWTSecret, u.ID, u.Email, u.Role, 12*time.Hour)
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   a.Cfg.AppEnv == "production",
		SameSite: "Lax",
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	return c.JSON(fiber.Map{"token": token, "user": fiber.Map{"id": u.ID, "email": u.Email, "role": u.Role}})
}

// POST /api/admin/logout
func (a *Admin) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{Name: "admin_token", Value: "", Expires: time.Now().Add(-time.Hour), Path: "/"})
	return c.JSON(fiber.Map{"ok": true})
}

// GET /api/admin/dashboard
func (a *Admin) Dashboard(c *fiber.Ctx) error {
	type counts struct {
		Total, Published, Hidden, Blocked, Categories, Tags, Reports int64
	}
	var ct counts
	a.DB.Model(&models.Post{}).Count(&ct.Total)
	a.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&ct.Published)
	a.DB.Model(&models.Post{}).Where("status = ?", "hidden").Count(&ct.Hidden)
	a.DB.Model(&models.Post{}).Where("status = ?", "blocked").Count(&ct.Blocked)
	a.DB.Model(&models.Category{}).Count(&ct.Categories)
	a.DB.Model(&models.Tag{}).Count(&ct.Tags)
	a.DB.Model(&models.Report{}).Where("status = ?", "open").Count(&ct.Reports)
	var lastLog models.ScrapeLog
	a.DB.Order("created_at DESC").First(&lastLog)
	return c.JSON(fiber.Map{"counts": ct, "last_scrape": lastLog})
}

// GET /api/admin/posts
func (a *Admin) ListPosts(c *fiber.Ctx) error {
	page, per := paginate(c)
	q := a.DB.Model(&models.Post{})
	if s := c.Query("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	if s := c.Query("safety_status"); s != "" {
		q = q.Where("safety_status = ?", s)
	}
	if s := c.Query("q"); s != "" {
		like := "%" + s + "%"
		q = q.Where("title ILIKE ? OR slug ILIKE ?", like, like)
	}
	var total int64
	q.Count(&total)
	var posts []models.Post
	q.Preload("Categories").Preload("Tags").
		Order("created_at DESC").Limit(per).Offset((page - 1) * per).Find(&posts)
	return c.JSON(fiber.Map{"data": posts, "page": page, "per_page": per, "total": total})
}

// PUT /api/admin/posts/:id
func (a *Admin) UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var p models.Post
	if err := a.DB.First(&p, id).Error; err != nil {
		return fiber.ErrNotFound
	}
	var in map[string]any
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	allowed := []string{"title", "excerpt", "content", "thumbnail_url", "video_embed_url", "status", "safety_status"}
	updates := map[string]any{}
	for _, k := range allowed {
		if v, ok := in[k]; ok {
			updates[k] = v
		}
	}
	if v, ok := updates["content"]; ok {
		if s, ok := v.(string); ok {
			updates["content"] = utils.SanitizeHTML(s)
		}
	}
	a.DB.Model(&p).Updates(updates)
	return c.JSON(fiber.Map{"data": p})
}

// DELETE /api/admin/posts/:id
func (a *Admin) DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := a.DB.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}

// PATCH /api/admin/posts/:id/status
func (a *Admin) UpdatePostStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var in struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	switch in.Status {
	case "published", "draft", "hidden", "blocked":
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid status")
	}
	if err := a.DB.Model(&models.Post{}).Where("id = ?", id).Update("status", in.Status).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}

// Categories CRUD
func (a *Admin) ListCategories(c *fiber.Ctx) error {
	var cats []models.Category
	a.DB.Order("name ASC").Find(&cats)
	return c.JSON(fiber.Map{"data": cats})
}
func (a *Admin) CreateCategory(c *fiber.Ctx) error {
	var in models.Category
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	if in.Slug == "" {
		in.Slug = utils.Slugify(in.Name)
	}
	if err := a.DB.Create(&in).Error; err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{"data": in})
}
func (a *Admin) UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var in models.Category
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	a.DB.Model(&models.Category{}).Where("id = ?", id).Updates(map[string]any{"name": in.Name, "slug": in.Slug})
	return c.JSON(fiber.Map{"ok": true})
}
func (a *Admin) DeleteCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	a.DB.Delete(&models.Category{}, id)
	return c.JSON(fiber.Map{"ok": true})
}

// Tags CRUD
func (a *Admin) ListTags(c *fiber.Ctx) error {
	var tags []models.Tag
	a.DB.Order("name ASC").Find(&tags)
	return c.JSON(fiber.Map{"data": tags})
}
func (a *Admin) CreateTag(c *fiber.Ctx) error {
	var in models.Tag
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	if in.Slug == "" {
		in.Slug = utils.Slugify(in.Name)
	}
	if err := a.DB.Create(&in).Error; err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{"data": in})
}
func (a *Admin) UpdateTag(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var in models.Tag
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	a.DB.Model(&models.Tag{}).Where("id = ?", id).Updates(map[string]any{"name": in.Name, "slug": in.Slug})
	return c.JSON(fiber.Map{"ok": true})
}
func (a *Admin) DeleteTag(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	a.DB.Delete(&models.Tag{}, id)
	return c.JSON(fiber.Map{"ok": true})
}

// Scraper
func (a *Admin) ScraperRun(c *fiber.Ctx) error {
	if a.Scheduler == nil {
		// Allow manual run via fresh scraper
		s := scraper.New(a.DB, a.Cfg)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()
			s.RunOnce(ctx)
		}()
		return c.JSON(fiber.Map{"ok": true, "started": true})
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		a.Scheduler.RunNow(ctx)
	}()
	return c.JSON(fiber.Map{"ok": true, "started": true})
}

func (a *Admin) ScraperLogs(c *fiber.Ctx) error {
	page, per := paginate(c)
	var total int64
	var logs []models.ScrapeLog
	a.DB.Model(&models.ScrapeLog{}).Count(&total)
	a.DB.Order("created_at DESC").Limit(per).Offset((page - 1) * per).Find(&logs)
	return c.JSON(fiber.Map{"data": logs, "page": page, "per_page": per, "total": total})
}

// Reports
func (a *Admin) ListReports(c *fiber.Ctx) error {
	page, per := paginate(c)
	var total int64
	var rs []models.Report
	q := a.DB.Model(&models.Report{})
	if s := c.Query("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	q.Count(&total)
	q.Order("created_at DESC").Limit(per).Offset((page - 1) * per).Find(&rs)
	return c.JSON(fiber.Map{"data": rs, "page": page, "per_page": per, "total": total})
}
func (a *Admin) UpdateReport(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var in struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	a.DB.Model(&models.Report{}).Where("id = ?", id).Update("status", in.Status)
	return c.JSON(fiber.Map{"ok": true})
}

// Settings
func (a *Admin) ListSettings(c *fiber.Ctx) error {
	var ss []models.Setting
	a.DB.Find(&ss)
	out := map[string]string{}
	for _, s := range ss {
		out[s.Key] = s.Value
	}
	return c.JSON(fiber.Map{"data": out})
}
func (a *Admin) UpdateSettings(c *fiber.Ctx) error {
	var in map[string]string
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	for k, v := range in {
		var s models.Setting
		if err := a.DB.Where("key = ?", k).First(&s).Error; err == nil {
			a.DB.Model(&s).Update("value", v)
		} else {
			a.DB.Create(&models.Setting{Key: k, Value: v})
		}
	}
	return c.JSON(fiber.Map{"ok": true})
}

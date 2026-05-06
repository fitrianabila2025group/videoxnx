package controllers

import (
	"strconv"
	"strings"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Public struct{ DB *gorm.DB }

func paginate(c *fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	per, _ := strconv.Atoi(c.Query("per_page", "20"))
	if per <= 0 || per > 100 {
		per = 20
	}
	return page, per
}

func publishedScope(q *gorm.DB) *gorm.DB {
	return q.Where("status = ? AND safety_status <> ?", "published", "blocked")
}

// GET /api/posts
func (p *Public) ListPosts(c *fiber.Ctx) error {
	page, per := paginate(c)
	var total int64
	var posts []models.Post

	q := publishedScope(p.DB.Model(&models.Post{}))
	q.Count(&total)
	if err := q.Preload("Categories").Preload("Tags").
		Order("COALESCE(published_at, scraped_at) DESC").
		Limit(per).Offset((page - 1) * per).Find(&posts).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": posts, "page": page, "per_page": per, "total": total})
}

// GET /api/posts/:slug
func (p *Public) GetPost(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var post models.Post
	if err := publishedScope(p.DB).Preload("Categories").Preload("Tags").
		Where("slug = ?", slug).First(&post).Error; err != nil {
		return fiber.ErrNotFound
	}
	p.DB.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	// related: same first category, exclude current
	var related []models.Post
	if len(post.Categories) > 0 {
		p.DB.Joins("JOIN post_categories pc ON pc.post_id = posts.id").
			Where("pc.category_id = ? AND posts.id <> ?", post.Categories[0].ID, post.ID).
			Where("posts.status = ?", "published").
			Limit(8).Find(&related)
	}

	// prev / next by published date
	var prev, next models.Post
	publishedScope(p.DB).Where("COALESCE(published_at, scraped_at) < COALESCE(?, ?)", post.PublishedAt, post.ScrapedAt).
		Order("COALESCE(published_at, scraped_at) DESC").First(&prev)
	publishedScope(p.DB).Where("COALESCE(published_at, scraped_at) > COALESCE(?, ?)", post.PublishedAt, post.ScrapedAt).
		Order("COALESCE(published_at, scraped_at) ASC").First(&next)

	return c.JSON(fiber.Map{
		"data":    post,
		"related": related,
		"prev":    nilIfEmpty(prev),
		"next":    nilIfEmpty(next),
	})
}

func nilIfEmpty(p models.Post) interface{} {
	if p.ID == 0 {
		return nil
	}
	return p
}

// GET /api/categories
func (p *Public) ListCategories(c *fiber.Ctx) error {
	var cats []models.Category
	p.DB.Order("name ASC").Find(&cats)
	return c.JSON(fiber.Map{"data": cats})
}

// GET /api/categories/:slug/posts
func (p *Public) PostsByCategory(c *fiber.Ctx) error {
	slug := c.Params("slug")
	page, per := paginate(c)
	var cat models.Category
	if err := p.DB.Where("slug = ?", slug).First(&cat).Error; err != nil {
		return fiber.ErrNotFound
	}
	var total int64
	var posts []models.Post
	q := publishedScope(p.DB.Model(&models.Post{})).
		Joins("JOIN post_categories pc ON pc.post_id = posts.id AND pc.category_id = ?", cat.ID)
	q.Count(&total)
	q.Preload("Categories").Preload("Tags").
		Order("COALESCE(published_at, scraped_at) DESC").
		Limit(per).Offset((page - 1) * per).Find(&posts)
	return c.JSON(fiber.Map{"data": posts, "category": cat, "page": page, "per_page": per, "total": total})
}

// GET /api/tags
func (p *Public) ListTags(c *fiber.Ctx) error {
	var tags []models.Tag
	p.DB.Order("name ASC").Find(&tags)
	return c.JSON(fiber.Map{"data": tags})
}

// GET /api/tags/:slug/posts
func (p *Public) PostsByTag(c *fiber.Ctx) error {
	slug := c.Params("slug")
	page, per := paginate(c)
	var tag models.Tag
	if err := p.DB.Where("slug = ?", slug).First(&tag).Error; err != nil {
		return fiber.ErrNotFound
	}
	var total int64
	var posts []models.Post
	q := publishedScope(p.DB.Model(&models.Post{})).
		Joins("JOIN post_tags pt ON pt.post_id = posts.id AND pt.tag_id = ?", tag.ID)
	q.Count(&total)
	q.Preload("Categories").Preload("Tags").
		Order("COALESCE(published_at, scraped_at) DESC").
		Limit(per).Offset((page - 1) * per).Find(&posts)
	return c.JSON(fiber.Map{"data": posts, "tag": tag, "page": page, "per_page": per, "total": total})
}

// GET /api/search?q=
func (p *Public) Search(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))
	page, per := paginate(c)
	if q == "" {
		return c.JSON(fiber.Map{"data": []any{}, "total": 0})
	}
	like := "%" + strings.ToLower(q) + "%"
	var total int64
	var posts []models.Post
	base := publishedScope(p.DB.Model(&models.Post{})).
		Where("LOWER(title) LIKE ? OR LOWER(excerpt) LIKE ?", like, like)
	base.Count(&total)
	base.Preload("Categories").Preload("Tags").
		Order("COALESCE(published_at, scraped_at) DESC").
		Limit(per).Offset((page - 1) * per).Find(&posts)
	return c.JSON(fiber.Map{"data": posts, "page": page, "per_page": per, "total": total, "q": q})
}

// GET /api/trending  (most viewed last 30 days fallback to all-time)
func (p *Public) Trending(c *fiber.Ctx) error {
	page, per := paginate(c)
	var posts []models.Post
	publishedScope(p.DB).Preload("Categories").Preload("Tags").
		Order("view_count DESC, COALESCE(published_at, scraped_at) DESC").
		Limit(per).Offset((page - 1) * per).Find(&posts)
	return c.JSON(fiber.Map{"data": posts, "page": page, "per_page": per})
}

// POST /api/reports
func (p *Public) CreateReport(c *fiber.Ctx) error {
	var in struct {
		PostID  uint   `json:"post_id"`
		Reason  string `json:"reason"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	if in.PostID == 0 || in.Reason == "" {
		return fiber.NewError(fiber.StatusBadRequest, "post_id and reason required")
	}
	r := models.Report{PostID: in.PostID, Reason: in.Reason, Email: in.Email, Message: in.Message, Status: "open"}
	if err := p.DB.Create(&r).Error; err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{"data": r})
}

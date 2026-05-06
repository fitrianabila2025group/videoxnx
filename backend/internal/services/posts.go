package services

import (
	"errors"
	"time"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostInput represents normalized scraped data ready to be persisted.
type PostInput struct {
	Title         string
	Slug          string
	Excerpt       string
	Content       string
	ThumbnailURL  string
	VideoEmbedURL string
	SourceURL     string
	SourceDomain  string
	PublishedAt   *time.Time
	Categories    []string
	Tags          []string
}

// UpsertPost saves or updates a post safely with safety filter applied.
func UpsertPost(db *gorm.DB, in PostInput) (*models.Post, bool, error) {
	if in.Title == "" || in.SourceURL == "" {
		return nil, false, errors.New("title and source_url required")
	}
	if in.Slug == "" {
		in.Slug = utils.Slugify(in.Title)
	}

	cleanContent := utils.SanitizeHTML(in.Content)
	cleanExcerpt := utils.StripHTML(in.Excerpt)
	tagsAndCats := append([]string{}, in.Categories...)
	tagsAndCats = append(tagsAndCats, in.Tags...)
	safety := ScanSafety(in.Title, cleanExcerpt, utils.StripHTML(cleanContent), tagsAndCats...)

	status := "published"
	if safety.Status == "blocked" {
		status = "blocked"
	} else if safety.Status == "review" {
		status = "hidden"
	}

	post := models.Post{
		Title:         in.Title,
		Slug:          uniqueSlug(db, in.Slug, 0),
		Excerpt:       cleanExcerpt,
		Content:       cleanContent,
		ThumbnailURL:  in.ThumbnailURL,
		VideoEmbedURL: in.VideoEmbedURL,
		SourceURL:     in.SourceURL,
		SourceDomain:  in.SourceDomain,
		Status:        status,
		IsAdult:       true,
		SafetyStatus:  safety.Status,
		SafetyReason:  safety.Reason,
		PublishedAt:   in.PublishedAt,
		ScrapedAt:     time.Now(),
	}

	// Upsert by source_url
	var existing models.Post
	err := db.Where("source_url = ?", in.SourceURL).First(&existing).Error
	created := false
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&post).Error; err != nil {
			return nil, false, err
		}
		created = true
		existing = post
	} else if err != nil {
		return nil, false, err
	} else {
		// Update mutable fields, do not overwrite admin status if hidden/blocked manually
		updates := map[string]interface{}{
			"title":           post.Title,
			"excerpt":         post.Excerpt,
			"content":         post.Content,
			"thumbnail_url":   post.ThumbnailURL,
			"video_embed_url": post.VideoEmbedURL,
			"source_domain":   post.SourceDomain,
			"published_at":    post.PublishedAt,
			"scraped_at":      post.ScrapedAt,
			"safety_status":   post.SafetyStatus,
			"safety_reason":   post.SafetyReason,
		}
		// Only update status if previously safe & not manually moderated.
		if existing.SafetyStatus == "safe" && safety.Status != "safe" {
			updates["status"] = status
		}
		if err := db.Model(&existing).Updates(updates).Error; err != nil {
			return nil, false, err
		}
	}

	// Categories / Tags
	cats := upsertCategories(db, in.Categories)
	tags := upsertTags(db, in.Tags)
	if len(cats) > 0 {
		_ = db.Model(&existing).Association("Categories").Replace(cats)
	}
	if len(tags) > 0 {
		_ = db.Model(&existing).Association("Tags").Replace(tags)
	}

	return &existing, created, nil
}

func uniqueSlug(db *gorm.DB, base string, ignoreID uint) string {
	slug := base
	i := 1
	for {
		var c int64
		q := db.Model(&models.Post{}).Where("slug = ?", slug)
		if ignoreID > 0 {
			q = q.Where("id <> ?", ignoreID)
		}
		q.Count(&c)
		if c == 0 {
			return slug
		}
		i++
		slug = base + "-" + itoa(i)
		if i > 1000 {
			return slug
		}
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	buf := []byte{}
	for i > 0 {
		buf = append([]byte{byte('0' + i%10)}, buf...)
		i /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}

func upsertCategories(db *gorm.DB, names []string) []models.Category {
	out := []models.Category{}
	for _, n := range dedupe(names) {
		if n == "" {
			continue
		}
		c := models.Category{Name: n, Slug: utils.Slugify(n)}
		db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).Create(&c)
		var f models.Category
		if err := db.Where("slug = ?", c.Slug).First(&f).Error; err == nil {
			out = append(out, f)
		}
	}
	return out
}

func upsertTags(db *gorm.DB, names []string) []models.Tag {
	out := []models.Tag{}
	for _, n := range dedupe(names) {
		if n == "" {
			continue
		}
		t := models.Tag{Name: n, Slug: utils.Slugify(n)}
		db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).Create(&t)
		var f models.Tag
		if err := db.Where("slug = ?", t.Slug).First(&f).Error; err == nil {
			out = append(out, f)
		}
	}
	return out
}

func dedupe(in []string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

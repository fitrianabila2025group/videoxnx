package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:191;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"size:32;default:admin" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Post struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Title         string     `gorm:"size:512;not null" json:"title"`
	Slug          string     `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Excerpt       string     `gorm:"type:text" json:"excerpt"`
	Content       string     `gorm:"type:text" json:"content"`
	ThumbnailURL  string     `gorm:"size:1024" json:"thumbnail_url"`
	VideoEmbedURL string     `gorm:"size:1024" json:"video_embed_url"`
	SourceURL     string     `gorm:"uniqueIndex;size:1024" json:"source_url"`
	SourceDomain  string     `gorm:"size:255;index" json:"source_domain"`
	Status        string     `gorm:"size:32;default:published;index" json:"status"` // published|draft|hidden|blocked
	IsAdult       bool       `gorm:"default:true" json:"is_adult"`
	SafetyStatus  string     `gorm:"size:32;default:safe;index" json:"safety_status"` // safe|review|blocked
	SafetyReason  string     `gorm:"size:512" json:"safety_reason,omitempty"`
	ViewCount     uint64     `gorm:"default:0" json:"view_count"`
	PublishedAt   *time.Time `json:"published_at"`
	ScrapedAt     time.Time  `json:"scraped_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Categories []Category `gorm:"many2many:post_categories" json:"categories,omitempty"`
	Tags       []Tag      `gorm:"many2many:post_tags" json:"tags,omitempty"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:191;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;size:191;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:191;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;size:191;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ScrapeLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SourceURL     string    `gorm:"size:1024" json:"source_url"`
	Status        string    `gorm:"size:32;index" json:"status"` // success|partial|failed
	Message       string    `gorm:"type:text" json:"message"`
	ScrapedCount  int       `json:"scraped_count"`
	FailedCount   int       `json:"failed_count"`
	DurationMs    int64     `json:"duration_ms"`
	CreatedAt     time.Time `json:"created_at"`
}

type Report struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"index" json:"post_id"`
	Reason    string    `gorm:"size:255" json:"reason"`
	Email     string    `gorm:"size:255" json:"email"`
	Message   string    `gorm:"type:text" json:"message"`
	Status    string    `gorm:"size:32;default:open;index" json:"status"` // open|reviewing|closed
	CreatedAt time.Time `json:"created_at"`
}

type Setting struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"uniqueIndex;size:191;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

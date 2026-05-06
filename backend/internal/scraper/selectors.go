package scraper

// Selectors are CSS selectors used as fallback when WP REST API is unavailable.
// Override these via env or settings if the source theme differs.
type Selectors struct {
	PostCard      string
	PostTitle     string
	PostLink      string
	PostThumbnail string
	PostExcerpt   string
	Content       string
	Iframe        string
	Category      string
	Tag           string
	Pagination    string
	NextPage      string
	PublishedDate string
}

// DefaultSelectors covers the most common WordPress video themes.
func DefaultSelectors() Selectors {
	return Selectors{
		PostCard:      "article, .post, .video-item, .item-video, .loop-video",
		PostTitle:     "h2.entry-title a, h3.entry-title a, .video-title a, h2 a, h3 a",
		PostLink:      "h2.entry-title a, h3.entry-title a, .video-title a, a.post-thumbnail, h2 a, h3 a",
		PostThumbnail: "img.wp-post-image, .post-thumbnail img, .video-thumb img, img",
		PostExcerpt:   ".entry-summary, .excerpt, .description",
		Content:       ".entry-content, .post-content, .video-description, article .content",
		Iframe:        "iframe, .video-embed iframe, .player iframe",
		Category:      ".cat-links a, .post-categories a, a[rel=category], a[rel='category tag']",
		Tag:           ".tag-links a, .post-tags a, a[rel=tag]",
		Pagination:    ".pagination a, .nav-links a, .page-numbers",
		NextPage:      "a.next, a.next.page-numbers, .nav-previous a",
		PublishedDate: "time.entry-date, time, meta[property='article:published_time']",
	}
}

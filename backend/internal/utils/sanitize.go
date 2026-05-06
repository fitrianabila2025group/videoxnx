package utils

import "github.com/microcosm-cc/bluemonday"

var (
	htmlSanitizer = func() *bluemonday.Policy {
		p := bluemonday.UGCPolicy()
		// Allow common embed iframes (only from known providers list will be enforced separately)
		p.AllowElements("iframe", "video", "source")
		p.AllowAttrs("src", "width", "height", "allow", "allowfullscreen", "frameborder", "scrolling").OnElements("iframe")
		p.AllowAttrs("src", "type", "controls", "poster", "preload").OnElements("video", "source")
		p.AllowURLSchemes("http", "https")
		return p
	}()
)

// SanitizeHTML cleans untrusted HTML, allowing safe embed tags.
func SanitizeHTML(in string) string {
	return htmlSanitizer.Sanitize(in)
}

// StripHTML returns plain text by removing all tags.
func StripHTML(in string) string {
	return bluemonday.StrictPolicy().Sanitize(in)
}

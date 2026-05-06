package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	slugRe      = regexp.MustCompile(`[^a-z0-9]+`)
	multiDashRe = regexp.MustCompile(`-+`)
)

// Slugify creates a URL-safe slug from a string. ASCII-only fallback.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRe.ReplaceAllString(s, "-")
	s = multiDashRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		b := make([]byte, 6)
		_, _ = rand.Read(b)
		return "post-" + hex.EncodeToString(b)
	}
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}

package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/services"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/utils"
)

// wpPost is a minimal subset of WP REST API post fields we use.
type wpPost struct {
	ID            int    `json:"id"`
	Date          string `json:"date_gmt"`
	Slug          string `json:"slug"`
	Link          string `json:"link"`
	Title         struct{ Rendered string } `json:"title"`
	Excerpt       struct{ Rendered string } `json:"excerpt"`
	Content       struct{ Rendered string } `json:"content"`
	FeaturedMedia int    `json:"featured_media"`
	Categories    []int  `json:"categories"`
	Tags          []int  `json:"tags"`
	Embedded      struct {
		FeaturedMedia []struct {
			SourceURL string `json:"source_url"`
		} `json:"wp:featuredmedia"`
		Term [][]struct {
			Name     string `json:"name"`
			Taxonomy string `json:"taxonomy"`
			Slug     string `json:"slug"`
		} `json:"wp:term"`
	} `json:"_embedded"`
}

// FetchViaWPRest tries to fetch posts from a WordPress REST endpoint.
// Returns (inputs, ok). ok=false means the API was unavailable / blocked / not WP.
func (s *Scraper) FetchViaWPRest(ctx context.Context, page int) ([]services.PostInput, bool) {
	endpoint := fmt.Sprintf("%s/wp-json/wp/v2/posts?_embed=1&per_page=20&page=%d", s.cfg.SourceBaseURL, page)
	resp, err := s.http.Get(ctx, endpoint)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, false
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return nil, false
	}
	// Some fetchers (e.g. headless Chromium) wrap raw JSON in <html><body><pre>...</pre>.
	// Strip that wrapper so json.Unmarshal works.
	body = unwrapJSON(body)
	var raw []wpPost
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, false
	}
	inputs := make([]services.PostInput, 0, len(raw))
	for _, p := range raw {
		in := services.PostInput{
			Title:        utils.StripHTML(p.Title.Rendered),
			Slug:         p.Slug,
			Excerpt:      utils.StripHTML(p.Excerpt.Rendered),
			Content:      p.Content.Rendered,
			SourceURL:    p.Link,
			SourceDomain: hostOf(p.Link),
		}
		if t, err := time.Parse(time.RFC3339, p.Date+"Z"); err == nil {
			in.PublishedAt = &t
		}
		if len(p.Embedded.FeaturedMedia) > 0 {
			in.ThumbnailURL = p.Embedded.FeaturedMedia[0].SourceURL
		}
		for _, group := range p.Embedded.Term {
			for _, term := range group {
				switch term.Taxonomy {
				case "category":
					in.Categories = append(in.Categories, term.Name)
				case "post_tag":
					in.Tags = append(in.Tags, term.Name)
				}
			}
		}
		// Try to find iframe in content
		in.VideoEmbedURL = extractIframeFromHTML(p.Content.Rendered)
		inputs = append(inputs, in)
	}
	return inputs, true
}

// FetchListPageHTML scrapes a listing/category/tag/page-N page and returns post detail URLs.
func (s *Scraper) FetchListPageHTML(ctx context.Context, pageURL string) ([]string, error) {
	resp, err := s.http.Get(ctx, pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	links := []string{}
	seen := map[string]struct{}{}
	doc.Find(s.sel.PostCard).Each(func(_ int, card *goquery.Selection) {
		a := card.Find(s.sel.PostLink).First()
		href, ok := a.Attr("href")
		if !ok {
			return
		}
		abs := absURL(pageURL, href)
		if _, dup := seen[abs]; dup {
			return
		}
		seen[abs] = struct{}{}
		links = append(links, abs)
	})
	return links, nil
}

// FetchPostHTML scrapes a single post page.
func (s *Scraper) FetchPostHTML(ctx context.Context, postURL string) (*services.PostInput, error) {
	resp, err := s.http.Get(ctx, postURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	in := &services.PostInput{
		SourceURL:    postURL,
		SourceDomain: hostOf(postURL),
	}
	in.Title = strings.TrimSpace(doc.Find("h1.entry-title, h1.post-title, h1").First().Text())
	if in.Title == "" {
		in.Title = strings.TrimSpace(doc.Find("title").Text())
	}
	in.Excerpt = strings.TrimSpace(doc.Find(s.sel.PostExcerpt).First().Text())
	contentHTML, _ := doc.Find(s.sel.Content).First().Html()
	in.Content = contentHTML

	// Thumbnail: og:image or first image
	if v, ok := doc.Find("meta[property='og:image']").Attr("content"); ok {
		in.ThumbnailURL = v
	} else {
		in.ThumbnailURL, _ = doc.Find(s.sel.PostThumbnail).First().Attr("src")
	}

	// Iframe / video — try real src then lazy-loaded data-src; also handle <video> and itemprop meta.
	in.VideoEmbedURL = extractVideoFromDoc(doc)

	// Categories
	doc.Find(s.sel.Category).Each(func(_ int, n *goquery.Selection) {
		t := strings.TrimSpace(n.Text())
		if t != "" {
			in.Categories = append(in.Categories, t)
		}
	})
	doc.Find(s.sel.Tag).Each(func(_ int, n *goquery.Selection) {
		t := strings.TrimSpace(n.Text())
		if t != "" {
			in.Tags = append(in.Tags, t)
		}
	})

	// Published date
	if v, ok := doc.Find("meta[property='article:published_time']").Attr("content"); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			in.PublishedAt = &t
		}
	} else if v, ok := doc.Find("time").First().Attr("datetime"); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			in.PublishedAt = &t
		}
	}

	if in.Title == "" {
		return nil, fmt.Errorf("no title found")
	}
	return in, nil
}

func extractIframeFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return extractVideoFromDoc(doc)
}

// extractVideoFromDoc looks at iframes (with lazy variants), <video src>/<source src>,
// and Schema.org itemprop="contentURL" / "embedURL" meta tags. Returns the first match.
func extractVideoFromDoc(doc *goquery.Document) string {
	// 1. iframe with src or lazy-loaded variants
	var found string
	doc.Find("iframe").EachWithBreak(func(_ int, n *goquery.Selection) bool {
		for _, attr := range []string{"src", "data-src", "data-litespeed-src", "data-lazy-src"} {
			if v, ok := n.Attr(attr); ok && strings.HasPrefix(v, "http") {
				found = v
				return false
			}
		}
		return true
	})
	if found != "" {
		return found
	}
	// 2. <video src> or nested <source src>
	doc.Find("video").EachWithBreak(func(_ int, n *goquery.Selection) bool {
		if v, ok := n.Attr("src"); ok && strings.HasPrefix(v, "http") {
			found = v
			return false
		}
		src, _ := n.Find("source").First().Attr("src")
		if strings.HasPrefix(src, "http") {
			found = src
			return false
		}
		return true
	})
	if found != "" {
		return found
	}
	// 3. itemprop=contentURL/embedURL meta
	for _, prop := range []string{"contentURL", "embedURL"} {
		if v, ok := doc.Find("meta[itemprop='" + prop + "']").First().Attr("content"); ok && strings.HasPrefix(v, "http") {
			return v
		}
	}
	return ""
}

func absURL(base, href string) string {
	b, err := url.Parse(base)
	if err != nil {
		return href
	}
	r, err := url.Parse(href)
	if err != nil {
		return href
	}
	return b.ResolveReference(r).String()
}

func hostOf(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Host
}

// satisfy unused import lint when http not used
var _ = http.StatusOK

// unwrapJSON returns the raw JSON bytes from a body. If the body has been
// wrapped by Chromium's JSON viewer (`<html>...<pre>{...}</pre>...</html>`),
// extract the contents of <pre>. Otherwise return as-is.
func unwrapJSON(b []byte) []byte {
	s := strings.TrimSpace(string(b))
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return []byte(s)
	}
	if i := strings.Index(s, "<pre"); i >= 0 {
		if start := strings.Index(s[i:], ">"); start >= 0 {
			rest := s[i+start+1:]
			if end := strings.Index(rest, "</pre>"); end >= 0 {
				inner := strings.TrimSpace(rest[:end])
				// HTML entity decode the few entities the json viewer emits
				inner = strings.NewReplacer("&quot;", "\"", "&amp;", "&", "&lt;", "<", "&gt;", ">").Replace(inner)
				return []byte(inner)
			}
		}
	}
	return b
}

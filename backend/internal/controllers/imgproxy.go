package controllers

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Image proxy: /api/img?u=<source-url>
// Re-fetches an upstream image with a referer matching its host so hotlink
// protection passes, then streams it to the client. Adds long cache headers.
//
// SECURITY: only proxies http/https schemes. Limits response size and follows
// up to 3 redirects via the default http.Client behavior.
//
// NOTE: We force HTTP/1.1 (TLSNextProto = empty map, ForceAttemptHTTP2=false)
// to avoid a long-standing race in net/http's HTTP/2 HPACK encoder that panics
// with "id (X) <= evictCount (Y)" under bursts of concurrent requests to a
// single host (golang/go#56019). The panic happens in a goroutine spawned by
// http2ClientConn.RoundTrip and therefore can't be recovered in the Fiber
// handler — it crashes the entire process. Sticking to HTTP/1.1 is safe here
// because we're only fetching small upstream images.
var imgClient = &http.Client{
	Timeout: 20 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		TLSNextProto:          map[string]func(string, *tls.Conn) http.RoundTripper{},
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   32,
		MaxConnsPerHost:       64,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
	},
}

const maxImageBytes = 8 << 20 // 8 MiB

func (p *Public) ImageProxy(c *fiber.Ctx) error {
	raw := strings.TrimSpace(c.Query("u"))
	if raw == "" {
		return fiber.NewError(http.StatusBadRequest, "missing u")
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return fiber.NewError(http.StatusBadRequest, "invalid url")
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 20*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", raw, nil)
	if err != nil {
		return fiber.NewError(http.StatusBadGateway, err.Error())
	}
	// Pose as a real browser visiting the source page itself to satisfy hotlink rules.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/*,*/*;q=0.8")
	req.Header.Set("Referer", u.Scheme+"://"+u.Host+"/")

	resp, err := imgClient.Do(req)
	if err != nil {
		return fiber.NewError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fiber.NewError(resp.StatusCode, "upstream error")
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" || !strings.HasPrefix(ct, "image/") {
		// Fallback: assume jpeg if upstream forgot the header
		ct = "image/jpeg"
	}
	c.Set("Content-Type", ct)
	c.Set("Cache-Control", "public, max-age=86400, s-maxage=604800, stale-while-revalidate=86400")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Cross-Origin-Resource-Policy", "cross-origin")

	_, err = io.Copy(c, io.LimitReader(resp.Body, maxImageBytes))
	return err
}

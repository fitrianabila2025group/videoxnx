package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/temoto/robotstxt"
)

var userAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36",
}

// HTTPClient handles polite HTTP requests with rate limiting, UA rotation, robots.txt.
// When FetcherURL is set, GET requests are proxied through the botasaurus fetcher
// microservice (which can solve Cloudflare challenges legitimately).
type HTTPClient struct {
	cli           *http.Client
	rateLimit     time.Duration
	defaultUA     string
	respectRobots bool
	robots        map[string]*robotstxt.RobotsData
	last          time.Time
	FetcherURL    string
}

func NewHTTPClient(rateLimitMs int, defaultUA string, respectRobots bool) *HTTPClient {
	return &HTTPClient{
		cli:           &http.Client{Timeout: 120 * time.Second},
		rateLimit:     time.Duration(rateLimitMs) * time.Millisecond,
		defaultUA:     defaultUA,
		respectRobots: respectRobots,
		robots:        map[string]*robotstxt.RobotsData{},
	}
}

func (c *HTTPClient) ua() string {
	if c.defaultUA != "" && rand.Intn(3) == 0 {
		return c.defaultUA
	}
	return userAgents[rand.Intn(len(userAgents))]
}

func (c *HTTPClient) waitRate() {
	if c.rateLimit <= 0 {
		return
	}
	d := time.Until(c.last.Add(c.rateLimit))
	if d > 0 {
		time.Sleep(d)
	}
	c.last = time.Now()
}

// Allowed checks robots.txt for the given URL.
func (c *HTTPClient) Allowed(ctx context.Context, target string) (bool, error) {
	if !c.respectRobots {
		return true, nil
	}
	u, err := url.Parse(target)
	if err != nil {
		return false, err
	}
	host := u.Scheme + "://" + u.Host
	r, ok := c.robots[host]
	if !ok {
		req, _ := http.NewRequestWithContext(ctx, "GET", host+"/robots.txt", nil)
		req.Header.Set("User-Agent", c.ua())
		resp, err := c.cli.Do(req)
		if err != nil {
			// allow if robots cannot be fetched
			c.robots[host] = nil
			return true, nil
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			c.robots[host] = nil
			return true, nil
		}
		r, err = robotstxt.FromBytes(body)
		if err != nil {
			c.robots[host] = nil
			return true, nil
		}
		c.robots[host] = r
	}
	if r == nil {
		return true, nil
	}
	return r.TestAgent(u.Path, c.defaultUA), nil
}

// Get performs a polite GET respecting robots.txt and rate limit.
// If FetcherURL is configured, requests are routed through the botasaurus microservice.
func (c *HTTPClient) Get(ctx context.Context, target string) (*http.Response, error) {
	if strings.TrimSpace(target) == "" {
		return nil, errors.New("empty url")
	}
	allowed, err := c.Allowed(ctx, target)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.New("disallowed by robots.txt: " + target)
	}
	c.waitRate()

	if c.FetcherURL != "" {
		return c.getViaFetcher(ctx, target)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.ua())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/json;q=0.9,*/*;q=0.5")
	req.Header.Set("Accept-Language", "en-US,en;q=0.7")
	return c.cli.Do(req)
}

// getViaFetcher proxies the request through the botasaurus fetcher microservice
// and synthesizes an *http.Response from the JSON payload it returns.
func (c *HTTPClient) getViaFetcher(ctx context.Context, target string) (*http.Response, error) {
	payload, _ := json.Marshal(map[string]any{"url": target, "wait": 4})
	req, err := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(c.FetcherURL, "/")+"/fetch", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 25<<20))
	if err != nil {
		return nil, err
	}
	var out struct {
		Status int    `json:"status"`
		HTML   string `json:"html"`
		URL    string `json:"url"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, errors.New("fetcher: invalid json: " + err.Error())
	}
	if out.Status >= 400 || out.Status == 0 {
		msg := out.Error
		if msg == "" {
			msg = "fetcher returned status " + http.StatusText(out.Status)
		}
		return nil, errors.New("fetcher: " + msg)
	}
	ct := "text/html; charset=utf-8"
	// Detect JSON responses (e.g. WP REST) so the parser can branch correctly.
	trimmed := strings.TrimSpace(out.HTML)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		ct = "application/json"
	}
	return &http.Response{
		StatusCode: out.Status,
		Status:     http.StatusText(out.Status),
		Body:       io.NopCloser(strings.NewReader(out.HTML)),
		Header:     http.Header{"Content-Type": []string{ct}},
		Request:    req,
	}, nil
}

package scraper

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/config"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/services"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
)

type Scraper struct {
	db   *gorm.DB
	cfg  *config.Config
	sel  Selectors
	http *HTTPClient
	mu   sync.Mutex
}

func New(db *gorm.DB, cfg *config.Config) *Scraper {
	hc := NewHTTPClient(cfg.ScraperRateLimitMs, cfg.ScraperUserAgent, cfg.ScraperRespectRobots)
	hc.FetcherURL = cfg.ScraperFetcherURL
	return &Scraper{
		db:   db,
		cfg:  cfg,
		sel:  DefaultSelectors(),
		http: hc,
	}
}

// RunOnce performs a full scrape pass. Safe to call manually or from scheduler.
func (s *Scraper) RunOnce(ctx context.Context) (saved int, failed int, err error) {
	if !s.mu.TryLock() {
		return 0, 0, fmt.Errorf("scraper already running")
	}
	defer s.mu.Unlock()

	start := time.Now()
	logEntry := &models.ScrapeLog{
		SourceURL: s.cfg.SourceBaseURL,
		Status:    "success",
	}
	defer func() {
		logEntry.ScrapedCount = saved
		logEntry.FailedCount = failed
		logEntry.DurationMs = time.Since(start).Milliseconds()
		if err != nil {
			logEntry.Status = "failed"
			logEntry.Message = err.Error()
		} else if failed > 0 {
			logEntry.Status = "partial"
		}
		s.db.Create(logEntry)
	}()

	// Try WP REST API first
	usedAPI := false
	consecutiveAllDup := 0
	for page := 1; page <= s.cfg.ScraperMaxPages; page++ {
		select {
		case <-ctx.Done():
			return saved, failed, ctx.Err()
		default:
		}
		inputs, ok := s.FetchViaWPRest(ctx, page)
		if !ok {
			if page == 1 {
				break
			}
			break
		}
		usedAPI = true
		if len(inputs) == 0 {
			break
		}
		pageNew := 0
		for _, in := range inputs {
			// WP REST often strips iframes from content.rendered; fetch the post page
			// to recover the embed if missing.
			if in.VideoEmbedURL == "" && in.SourceURL != "" {
				if html, err := s.FetchPostHTML(ctx, in.SourceURL); err == nil {
					if html.VideoEmbedURL != "" {
						in.VideoEmbedURL = html.VideoEmbedURL
					}
					if in.ThumbnailURL == "" {
						in.ThumbnailURL = html.ThumbnailURL
					}
				}
			}
			_, created, e := services.UpsertPost(s.db, in)
			if e != nil {
				failed++
				log.Printf("scraper upsert failed: %v", e)
			} else {
				saved++
				if created {
					pageNew++
				}
			}
		}
		log.Printf("scraper WP REST page=%d items=%d new=%d total_saved=%d", page, len(inputs), pageNew, saved)

		// Stop early once we've seen 3 consecutive pages with zero NEW posts
		// (steady-state: the source has nothing new to give us). Full backfill
		// happens automatically because every page that still introduces new
		// posts resets this counter.
		if pageNew == 0 {
			consecutiveAllDup++
			if consecutiveAllDup >= 3 {
				log.Printf("scraper: 3 consecutive pages without new posts, stopping at page=%d", page)
				break
			}
		} else {
			consecutiveAllDup = 0
		}
	}

	if usedAPI {
		return saved, failed, nil
	}

	// Fallback HTML scraping
	pageURL := s.cfg.SourceBaseURL
	consecutiveAllDup = 0
	for page := 1; page <= s.cfg.ScraperMaxPages; page++ {
		select {
		case <-ctx.Done():
			return saved, failed, ctx.Err()
		default:
		}
		links, e := s.FetchListPageHTML(ctx, pageURL)
		if e != nil {
			failed++
			log.Printf("list page failed (%s): %v", pageURL, e)
			break
		}
		pageNew := 0
		for _, link := range links {
			in, e := s.FetchPostHTML(ctx, link)
			if e != nil {
				failed++
				log.Printf("post fetch failed (%s): %v", link, e)
				continue
			}
			_, created, e := services.UpsertPost(s.db, *in)
			if e != nil {
				failed++
				log.Printf("upsert failed: %v", e)
			} else {
				saved++
				if created {
					pageNew++
				}
			}
		}
		log.Printf("scraper HTML page=%d items=%d new=%d total_saved=%d", page, len(links), pageNew, saved)
		if pageNew == 0 {
			consecutiveAllDup++
			if consecutiveAllDup >= 3 {
				log.Printf("scraper: 3 consecutive pages without new posts, stopping at page=%d", page)
				break
			}
		} else {
			consecutiveAllDup = 0
		}
		// pagination: WordPress style /page/N/
		pageURL = fmt.Sprintf("%s/page/%d/", s.cfg.SourceBaseURL, page+1)
	}

	return saved, failed, nil
}

// Scheduler wraps a Scraper with a cron schedule.
type Scheduler struct {
	s   *Scraper
	cfg *config.Config
	sch gocron.Scheduler
}

func NewScheduler(db *gorm.DB, cfg *config.Config) *Scheduler {
	sc, _ := gocron.NewScheduler()
	return &Scheduler{s: New(db, cfg), cfg: cfg, sch: sc}
}

func (s *Scheduler) Start() {
	interval := time.Duration(s.cfg.ScraperIntervalMin) * time.Minute
	if interval <= 0 {
		interval = time.Hour
	}
	_, err := s.sch.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
			defer cancel()
			if saved, failed, err := s.s.RunOnce(ctx); err != nil {
				log.Printf("scheduled scrape error: %v", err)
			} else {
				log.Printf("scheduled scrape: saved=%d failed=%d", saved, failed)
			}
		}),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		log.Printf("scheduler job error: %v", err)
		return
	}
	s.sch.Start()
}

func (s *Scheduler) Stop() {
	_ = s.sch.Shutdown()
}

func (s *Scheduler) RunNow(ctx context.Context) (int, int, error) {
	return s.s.RunOnce(ctx)
}

func (s *Scheduler) Scraper() *Scraper { return s.s }

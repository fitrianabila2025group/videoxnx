package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv               string
	AppURL               string
	Port                 string
	DatabaseURL          string
	RedisURL             string
	JWTSecret            string
	AdminEmail           string
	AdminPassword        string
	SourceBaseURL        string
	ScraperEnabled       bool
	ScraperIntervalMin   int
	ScraperMaxPages      int
	ScraperRateLimitMs   int
	ScraperRespectRobots bool
	ScraperUserAgent     string
	ScraperFetcherURL    string
	AgeGateEnabled       bool
	DMCAEmail            string
	SiteName             string
	MetaTitle            string
	MetaDescription      string
	CORSAllowedOrigins   []string
}

func Load() *Config {
	return &Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		AppURL:               getEnv("APP_URL", "http://localhost:8080"),
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		RedisURL:             getEnv("REDIS_URL", ""),
		JWTSecret:            getEnv("JWT_SECRET", "insecure-dev-secret"),
		AdminEmail:           getEnv("ADMIN_EMAIL", "mpratamagpt@gmail.com"),
		AdminPassword:        getEnv("ADMIN_PASSWORD", "Anonymous263"),
		SourceBaseURL:        strings.TrimRight(getEnv("SOURCE_BASE_URL", "https://indoxvx.cam"), "/"),
		ScraperEnabled:       getEnvBool("SCRAPER_ENABLED", true),
		ScraperIntervalMin:   getEnvInt("SCRAPER_INTERVAL_MINUTES", 60),
		ScraperMaxPages:      getEnvInt("SCRAPER_MAX_PAGES", 20),
		ScraperRateLimitMs:   getEnvInt("SCRAPER_RATE_LIMIT_MS", 1500),
		ScraperRespectRobots: getEnvBool("SCRAPER_RESPECT_ROBOTS", true),
		ScraperUserAgent:     getEnv("SCRAPER_USER_AGENT", "VideoxnxAggregatorBot/1.0"),
		ScraperFetcherURL:    getEnv("SCRAPER_FETCHER_URL", ""),
		AgeGateEnabled:       getEnvBool("AGE_GATE_ENABLED", true),
		DMCAEmail:            getEnv("DMCA_EMAIL", ""),
		SiteName:             getEnv("SITE_NAME", "VideoXNX"),
		MetaTitle:            getEnv("META_TITLE", "VideoXNX"),
		MetaDescription:      getEnv("META_DESCRIPTION", ""),
		CORSAllowedOrigins:   splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "*")),
	}
}

func getEnv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		return v
	}
	return def
}
func getEnvInt(k string, def int) int {
	if v, ok := os.LookupEnv(k); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
func getEnvBool(k string, def bool) bool {
	if v, ok := os.LookupEnv(k); ok {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

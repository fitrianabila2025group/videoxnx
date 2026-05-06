package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/config"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/database"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/routes"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/scraper"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("db migrate: %v", err)
	}
	if err := services.EnsureAdminUser(db, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("admin seed: %v", err)
	}

	app := routes.NewApp(db, cfg)

	// Scraper scheduler
	var sched *scraper.Scheduler
	if cfg.ScraperEnabled {
		sched = scraper.NewScheduler(db, cfg)
		sched.Start()
	}

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if sched != nil {
			sched.Stop()
		}
		_ = app.ShutdownWithContext(ctx)
	}()

	log.Printf("API listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

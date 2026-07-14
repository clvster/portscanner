package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"portscanner/config"
	"portscanner/db"
	"portscanner/services/notifier"
	"portscanner/services/scanner"
	"portscanner/services/scheduler"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	database, err := db.New(ctx, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(ctx); err != nil {
		log.Fatalf("db migrate: %v", err)
	}
	log.Println("db connected and migrated")

	sc := scanner.New(cfg.Scan.Rate, cfg.Scan.SourceIP)
	notify := notifier.NewTelegram(cfg.Telegram.Token, cfg.Telegram.ChatID)
	sched := scheduler.New(sc, database, notify, cfg.Scan.Targets, cfg.Scan.Ports)

	log.Printf("targets=%v ports=%s rate=%d", cfg.Scan.Targets, cfg.Scan.Ports, cfg.Scan.Rate)

	if cfg.Scan.Once {
		if err := sched.RunOnce(ctx); err != nil {
			log.Fatalf("scan: %v", err)
		}
		return
	}

	log.Printf("running scheduled scans every %s", cfg.Scan.Interval)
	sched.Run(ctx, cfg.Scan.Interval)
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"portscanner/config"
	"portscanner/db"
	"portscanner/services/notifier"
	"portscanner/services/scanner"
	"portscanner/types"
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

	notify := notifier.NewTelegram(cfg.Telegram.Token, cfg.Telegram.ChatID)

	log.Printf("targets=%v ports=%s rate=%d", cfg.Scan.Targets, cfg.Scan.Ports, cfg.Scan.Rate)

	sc := scanner.New(cfg.Scan.Rate)

	log.Println("starting masscan...")
	ports, err := sc.Scan(ctx, cfg.Scan.Targets, cfg.Scan.Ports)
	if err != nil {
		log.Fatalf("scan: %v", err)
	}

	if len(ports) == 0 {
		log.Println("no open ports found")
		os.Exit(0)
	}

	log.Printf("found %d open port(s):", len(ports))
	var newPorts []types.OpenPort
	for _, p := range ports {
		isNew, err := database.UpsertPort(ctx, p)
		if err != nil {
			log.Printf("  upsert %s: %v", p.Key(), err)
			continue
		}
		state := "seen"
		if isNew {
			state = "NEW"
			newPorts = append(newPorts, p)
		}
		if p.Service != "" {
			log.Printf("  [%s] %s:%d/%s service=%s banner=%q", state, p.IP, p.Port, p.Proto, p.Service, p.Banner)
		} else {
			log.Printf("  [%s] %s:%d/%s", state, p.IP, p.Port, p.Proto)
		}
	}

	log.Printf("scan complete: %d new, %d total", len(newPorts), len(ports))

	if len(newPorts) > 0 {
		if err := notify.Notify(ctx, newPorts); err != nil {
			log.Printf("notify: %v", err)
		} else {
			log.Printf("notified about %d new port(s)", len(newPorts))
		}
	}
}

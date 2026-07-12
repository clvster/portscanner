package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"portscanner/config"
	"portscanner/services/scanner"
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
	for _, p := range ports {
		if p.Service != "" {
			log.Printf("  %s:%d/%s  service=%s  banner=%q", p.IP, p.Port, p.Proto, p.Service, p.Banner)
		} else {
			log.Printf("  %s:%d/%s", p.IP, p.Port, p.Proto)
		}
	}
}

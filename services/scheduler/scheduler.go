package scheduler

import (
	"context"
	"log"
	"time"

	"portscanner/db"
	"portscanner/services/notifier"
	"portscanner/services/scanner"
	"portscanner/types"
)

type Scheduler struct {
	scanner  *scanner.Masscan
	db       *db.DB
	notifier notifier.Notifier
	targets  []string
	ports    string
}

func New(sc *scanner.Masscan, database *db.DB, n notifier.Notifier, targets []string, ports string) *Scheduler {
	return &Scheduler{
		scanner:  sc,
		db:       database,
		notifier: n,
		targets:  targets,
		ports:    ports,
	}
}

func (s *Scheduler) RunOnce(ctx context.Context) error {
	log.Println("scan cycle started")

	ports, err := s.scanner.Scan(ctx, s.targets, s.ports)
	if err != nil {
		return err
	}

	if len(ports) == 0 {
		log.Println("no open ports found")
		return nil
	}

	log.Printf("found %d open port(s):", len(ports))
	var newPorts []types.OpenPort

	for _, p := range ports {
		isNew, err := s.db.UpsertPort(ctx, p)
		if err != nil {
			log.Printf("    upsert %s: %v", p.Key(), err)
		}

		state := "seen"

		if isNew {
			state = "NEW"
			newPorts = append(newPorts, p)
		}

		if p.Service != "" {
			log.Printf("    [%s] %s:%d/%s service=%s banner=%q", state, p.IP, p.Port, p.Proto, p.Service, p.Banner)
		} else {
			log.Printf("    [%s] %s:%d/%s", state, p.IP, p.Port, p.Proto)
		}
	}

	log.Printf("scan complete: %d new, %d total", len(newPorts), len(ports))
	if len(newPorts) > 0 {
		if err := s.notifier.Notify(ctx, newPorts); err != nil {
			log.Printf("notify: %v", err)
		} else {
			log.Printf("notified about %d new port(s)", len(newPorts))
		}
	}

	return nil
}

func (s *Scheduler) Run(ctx context.Context, interval time.Duration) {
	if err := s.RunOnce(ctx); err != nil {
		log.Printf("scan cycle: %v", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutdown signal received, exiting")
			return
		case <-ticker.C:
			if err := s.RunOnce(ctx); err != nil {
				log.Printf("scan cycle: %v", err)
			}
		}
	}
}

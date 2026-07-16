package db

import (
	"context"
	"fmt"

	"portscanner/types"
)

func (d *DB) UpsertPort(ctx context.Context, p types.OpenPort) (isNew bool, err error) {
	const q = `
		INSERT INTO open_ports (ip, port, proto, service, banner, product, version, cpe, cves)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (ip, port, proto)
		DO UPDATE SET
			last_seen = now(),
			service = EXCLUDED.service,
			banner = EXCLUDED.banner,
			product = EXCLUDED.product,
			version = EXCLUDED.version,
			cpe = EXCLUDED.cpe,
			cves = EXCLUDED.cves
		RETURNING (xmax = 0) AS inserted;
	`

	cves := p.CVEs
	if cves == nil {
		cves = []string{}
	}

	err = d.pool.QueryRow(ctx, q, p.IP, p.Port, p.Proto, p.Service, p.Banner, p.Product, p.Version, p.CPE, cves).Scan(&isNew)
	if err != nil {
		return false, fmt.Errorf("upsert port %s: %w", p.Key(), err)
	}

	return isNew, nil
}

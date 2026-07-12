package db

import (
	"context"
	"fmt"

	"portscanner/types"
)

func (d *DB) UpsertPort(ctx context.Context, p types.OpenPort) (isNew bool, err error) {
	const q = `
		INSERT INTO open_ports (ip, port, proto, service, banner)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ip, port, proto)
		DO UPDATE SET
			last_seen = now(),
			service   = EXCLUDED.service,
			banner    = EXCLUDED.banner
		RETURNING (xmax = 0) AS inserted;
	`

	err = d.pool.QueryRow(ctx, q, p.IP, p.Port, p.Proto, p.Service, p.Banner).Scan(&isNew)
	if err != nil {
		return false, fmt.Errorf("upsert port %s: %w", p.Key(), err)
	}

	return isNew, nil
}

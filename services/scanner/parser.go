package scanner

import (
	"encoding/json"
	"strings"

	"portscanner/types"
)

type rawEntry struct {
	IP        string    `json:"ip"`
	Timestamp string    `json:"timestamp"`
	Ports     []rawPort `json:"ports"`
}

type rawPort struct {
	Port    int         `json:"port"`
	Proto   string      `json:"proto"`
	Status  string      `json:"status"`
	Service *rawService `json:"service"`
}

type rawService struct {
	Name   string `json:"name"`
	Banner string `json:"banner"`
}

func parseMasscanJSON(data []byte) ([]types.OpenPort, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "[]" {
		return nil, nil
	}

	var entries []rawEntry
	if err := json.Unmarshal([]byte(trimmed), &entries); err != nil {
		return nil, err
	}

	merged := make(map[string]types.OpenPort)
	order := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.IP == "" {
			continue
		}
		for _, p := range e.Ports {
			if p.Port == 0 && p.Proto == "" {
				continue
			}
			proto := p.Proto
			if proto == "" {
				proto = "tcp"
			}

			op := types.OpenPort{
				IP:    e.IP,
				Port:  p.Port,
				Proto: proto,
			}
			if p.Service != nil {
				op.Service = p.Service.Name
				op.Banner = p.Service.Banner
			}

			key := op.Key()
			if existing, ok := merged[key]; ok {
				if op.Service != "" {
					existing.Service = op.Service
				}
				if op.Banner != "" {
					existing.Banner = op.Banner
				}
				merged[key] = existing
			} else {
				merged[key] = op
				order = append(order, key)
			}
		}
	}

	result := make([]types.OpenPort, 0, len(order))
	for _, key := range order {
		result = append(result, merged[key])
	}

	return result, nil
}

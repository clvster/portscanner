package enricher

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"portscanner/types"
)

type Nmap struct {
	binary string
}

func NewNmap() *Nmap {
	return &Nmap{binary: "nmap"}
}

func (n *Nmap) Enrich(ctx context.Context, ports []types.OpenPort) []types.OpenPort {
	if len(ports) == 0 {
		return ports
	}

	byIP := make(map[string][]int)
	for _, p := range ports {
		byIP[p.IP] = append(byIP[p.IP], p.Port)
	}

	found := make(map[string]enriched)
	for ip, portList := range byIP {
		results, err := n.scanHost(ctx, ip, portList)
		if err != nil {
			continue
		}
		for _, e := range results {
			found[e.IP+":"+strconv.Itoa(e.Port)+"/"+e.Proto] = e
		}
	}

	out := make([]types.OpenPort, len(ports))
	for i, p := range ports {
		out[i] = p
		if e, ok := found[p.Key()]; ok {
			if p.Service == "" {
				out[i].Service = e.Service
			}
			out[i].Product = e.Product
			out[i].Version = e.Version
			out[i].CPE = e.CPE
		}
	}

	return out
}

func (n *Nmap) scanHost(ctx context.Context, ip string, ports []int) ([]enriched, error) {
	portArgs := make([]string, len(ports))
	for i, p := range ports {
		portArgs[i] = strconv.Itoa(p)
	}

	args := []string{
		"-sV",
		"-Pn",
		"-p", strings.Join(portArgs, ","),
		"-oX", "-",
		ip,
	}

	cmd := exec.CommandContext(ctx, n.binary, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nmap %s: %w", ip, err)
	}

	return parseNmapXML(out)
}

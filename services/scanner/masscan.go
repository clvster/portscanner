package scanner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"portscanner/types"
)

type Masscan struct {
	binary   string
	rate     int
	sourceIP string
}

func New(rate int, sourceIP string) *Masscan {
	return &Masscan{
		binary:   "masscan",
		rate:     rate,
		sourceIP: sourceIP,
	}
}

func (m *Masscan) Scan(ctx context.Context, targets []string, ports string) ([]types.OpenPort, error) {
	outFile, err := os.CreateTemp("", "masscan-*.json")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	outPath := outFile.Name()
	outFile.Close()
	defer os.Remove(outPath)

	args := []string{
		"-p", ports,
		"--rate", strconv.Itoa(m.rate),
		"--banners",
		"--open-only",
		"-oJ", outPath,
	}

	if m.sourceIP != "" {
		args = append(args, "--source-ip", m.sourceIP)
	}

	args = append(args, targets...)

	cmd := exec.CommandContext(ctx, m.binary, args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("masscan run: %w (stderr: %s)", err, stderr.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read output: %w", err)
	}

	return parseMasscanJSON(data)
}

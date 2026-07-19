package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"portscanner/types"
)

type Telegram struct {
	token  string
	chatID int64
	client *http.Client
}

func NewTelegram(token string, chatID int64) *Telegram {
	return &Telegram{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (t *Telegram) Notify(ctx context.Context, ports []types.OpenPort) error {
	if len(ports) == 0 {
		return nil
	}

	body, err := json.Marshal(map[string]any{
		"chat_id":    t.chatID,
		"text":       formatMessage(ports),
		"parse_mode": "HTML",
	})
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request to telegram api: %s", extractNetError(err))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api status %d: %s", resp.StatusCode, string(data))
	}

	return nil
}

func extractNetError(err error) string {
	msg := err.Error()
	if i := strings.Index(msg, "\": "); i != -1 {
		return msg[i+3:]
	}
	return "request failed"
}

func formatMessage(ports []types.OpenPort) string {
	var b strings.Builder
	fmt.Fprintf(&b, "<b>New open ports: %d</b>\n\n", len(ports))
	for _, p := range ports {
		fmt.Fprintf(&b, "• <code>%s:%d/%s</code>", p.IP, p.Port, p.Proto)
		switch {
		case p.Product != "":
			fmt.Fprintf(&b, " — %s", p.Product)
			if p.Version != "" {
				fmt.Fprintf(&b, " %s", p.Version)
			}
		case p.Service != "":
			fmt.Fprintf(&b, " — %s", p.Service)
		}
		if len(p.CVEs) > 0 {
			fmt.Fprintf(&b, "\n  %s", strings.Join(p.CVEs, ", "))
		}
		b.WriteString("\n")
	}
	return b.String()
}

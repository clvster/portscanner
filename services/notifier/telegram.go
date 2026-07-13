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
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *Telegram) Notify(ctx context.Context, ports []types.OpenPort) error {
	if len(ports) == 0 {
		return nil
	}

	text := formatMessage(ports)

	body, err := json.Marshal(map[string]any{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "HTML",
	})
	if err != nil {
		return fmt.Errorf("marshal message %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api status %d: %s", resp.StatusCode, string(data))
	}

	return nil
}

func formatMessage(ports []types.OpenPort) string {
	var b strings.Builder
	fmt.Fprintf(&b, "<b>New open ports: %d</b>\n\n", len(ports))

	for _, p := range ports {
		fmt.Fprintf(&b, "• <code>%s:%d/%s</code>", p.IP, p.Port, p.Proto)

		if p.Service != "" {
			fmt.Fprintf(&b, " — %s", p.Service)
		}
		b.WriteString("\n")
	}

	return b.String()
}

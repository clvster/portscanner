package notifier

import (
	"strings"
	"testing"

	"portscanner/types"
)

func TestFormatMessage(t *testing.T) {
	ports := []types.OpenPort{
		{IP: "172.18.0.2", Port: 80, Proto: "tcp", Service: "http"},
		{IP: "172.18.0.4", Port: 6379, Proto: "tcp"},
	}

	msg := formatMessage(ports)

	if !strings.Contains(msg, "New open ports: 2") {
		t.Errorf("missing count header: %q", msg)
	}
	if !strings.Contains(msg, "172.18.0.2:80/tcp") {
		t.Errorf("missing port 80 entry: %q", msg)
	}
	if !strings.Contains(msg, "— http") {
		t.Errorf("missing service label for port 80: %q", msg)
	}
	if !strings.Contains(msg, "172.18.0.4:6379/tcp") {
		t.Errorf("missing port 6379 entry: %q", msg)
	}
}

package notifier

import (
	"context"

	"portscanner/types"
)

type Notifier interface {
	Notify(ctx context.Context, ports []types.OpenPort) error
}

package ssh

import (
	"context"
	"time"
)

// contextTimeout is split out so the test of TestAuth can stub it later.
func contextTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

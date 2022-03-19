package retry

import (
	"context"
	"time"
)

type retryableFunc func() error
type checkRetryable func(err error) bool
type option func(*backoff)

func Do(ctx context.Context, fn retryableFunc, ops ...option) error {
	b := defaultBackoff()
	for _, op := range ops {
		op(b)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var err error
	var next time.Duration
	for {
		if err = fn(); err == nil {
			return nil
		}

		if !b.checkRetryable(err) {
			return err
		}

		if next = b.next(); next == stop {
			return err
		}

		t := time.NewTimer(next)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			// In case ctx.Done and t.C happen completely at the same time
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
	}
}

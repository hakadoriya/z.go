package contextz

import (
	"context"
	"fmt"
)

func CheckContext(ctx context.Context) error {
	select {
	default:
		return nil
	case <-ctx.Done():
		if err := context.Cause(ctx); err != nil {
			return fmt.Errorf("ctx.Err() = %w: context.Cause: %w", ctx.Err(), err)
		}
		return fmt.Errorf("ctx.Err: %w", ctx.Err())
	}
}

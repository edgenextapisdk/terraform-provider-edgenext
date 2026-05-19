package elb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// elbOctaviaImmutableRetryable is true when Octavia rejects the call because a load balancer or target group is still applying a prior async change (409 immutable, etc.).
func elbOctaviaImmutableRetryable(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "immutable") ||
		strings.Contains(s, "cannot be updated") ||
		strings.Contains(s, "code: 409") ||
		strings.Contains(s, " 409")
}

// elbPOSTWithOctaviaImmutableRetry performs an ELB POST and retries on transient Octavia immutable/conflict errors.
func elbPOSTWithOctaviaImmutableRetry(ctx context.Context, elbClient *connectivity.ELBClient, path string, body interface{}, resp *map[string]interface{}) error {
	const maxAttempts = 20
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(2+attempt) * time.Second
			if backoff > 8*time.Second {
				backoff = 8 * time.Second
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}
		err := elbClient.Post(ctx, path, body, resp)
		if err == nil {
			return nil
		}
		lastErr = err
		if !elbOctaviaImmutableRetryable(err) {
			return err
		}
	}
	if lastErr == nil {
		return fmt.Errorf("POST %s: exhausted retries", path)
	}
	return fmt.Errorf("POST %s: exhausted retries: %w", path, lastErr)
}

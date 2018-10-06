// Package retry provides retry functions that can be used
// when talking to unreliable services.
package retry

import (
	"math"
	"time"
)

// BackoffRetryN attempts at most n retries, exponentially backing off from minDelay
// but never waiting more than maxDelay before the next retry. The function f is called
// for each retry and should be idempotent. nil is returned on success. On error the result
// of the last call to f is returned.
func BackoffRetryN(n int, minDelay, maxDelay time.Duration, f func() error) error {
	delay := minDelay
	for i := 1; i < n; i++ {
		if err := f(); err == nil {
			return nil
		}
		time.Sleep(delay)
		delay = time.Duration(math.Min(float64(minDelay)*math.Pow(2, float64(i)), float64(maxDelay)))
	}

	// On the nth time, try without a delay following.
	return f()
}

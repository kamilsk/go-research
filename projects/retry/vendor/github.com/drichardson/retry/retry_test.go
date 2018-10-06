package retry

import (
	"errors"
	"testing"
	"time"
)

func TestBackoffRetryN(t *testing.T) {
	counter := 0
	ErrTest := errors.New("ErrTest")
	before := time.Now()
	err := BackoffRetryN(5, 1*time.Millisecond, 1*time.Hour, func() error {
		counter++
		return ErrTest
	})
	after := time.Now()
	if err != ErrTest {
		t.Errorf("expected ErrTest but got %v", err)
	}
	if counter != 5 {
		t.Errorf("counter 5 != %v", counter)
	}
	// should have waited a total of 1ms after first execution, 2ms after
	// second, 4ms after third, 8ms after fourth, and 0ms after fifth (final).
	// 1ms+2ms+4ms+8ms=15ms
	expectedMinimumDuration := 15 * time.Millisecond
	diff := after.Sub(before)
	if diff < expectedMinimumDuration {
		t.Errorf("Expected delay of at least %v but got %v", expectedMinimumDuration, diff)
	}
}

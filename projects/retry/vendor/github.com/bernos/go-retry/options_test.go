package retry

import (
	"fmt"
	"testing"
	"time"
)

func TestBaseDelay(t *testing.T) {
	r := &Options{}

	BaseDelay(time.Second)(r)

	if r.BaseDelay != time.Second {
		t.Errorf("Want %s, got %s", time.Second, r.BaseDelay)
	}
}

func TestForever(t *testing.T) {
	r := &Options{}

	Forever()(r)

	if r.MaxRetries != Infinity {
		t.Errorf("Want %d, got %d", Infinity, r.MaxRetries)
	}
}

func TestMaxDelay(t *testing.T) {
	r := &Options{}

	MaxDelay(time.Hour)(r)

	if r.MaxDelay != time.Hour {
		t.Errorf("Want %s, got %s", time.Hour, r.MaxDelay)
	}
}

func TestExponentialBackoff(t *testing.T) {
	r := &Options{}

	BinaryExponentialBackoff()(r)

	if r.CalculateDelay(5, time.Millisecond, time.Minute) != calculateBinaryExponentialDelay(5, time.Millisecond, time.Minute) {
		t.Errorf("Want calculateDelayBinary")
	}
}

func TestWhile(t *testing.T) {
	count := 0
	max := 3

	fn := func() (interface{}, error) {
		count++
		return nil, fmt.Errorf("Test retry")
	}

	r := Retry(fn, While(func(error) bool {
		return count < max
	}))

	_, err := r()

	if err == nil {
		t.Errorf("Expected an error")
	}

	if count != max {
		t.Errorf("Want %d, got %d", max, count)
	}
}

func TestLog(t *testing.T) {
	count := 0
	messageCount := 0

	fn := func() (interface{}, error) {
		if count == 10 {
			return "done", nil
		}
		count++
		return nil, fmt.Errorf("Test retry")
	}

	r := Retry(fn, Log(func(format string, v ...interface{}) {
		messageCount++
	}))

	v, err := r()

	if v != "done" {
		t.Errorf("Want %s, got %s", "done", v)
	}

	if err != nil {
		t.Errorf("Unexpected error. %s", err.Error())
	}

	if messageCount < count {
		t.Errorf("Want %d messages, got %d", count, messageCount)
	}
}

package retry

import (
	"fmt"
	"math"
	"testing"
	"testing/quick"
	"time"
)

func TestCalculateDelayBinary(t *testing.T) {
	maxInt := int((^uint(0)) >> 1)

	tests := []struct {
		min       int
		max       int
		baseDelay time.Duration
		maxDelay  time.Duration
	}{
		{0, 100, time.Millisecond, time.Minute},
		{maxInt - 10, maxInt + 10, time.Nanosecond, time.Minute},
	}

	for _, test := range tests {
		last := time.Duration(0)

		for i := test.min; i < test.max; i++ {
			d := calculateBinaryExponentialDelay(uint(i), test.baseDelay, test.maxDelay)
			if d < 0 || d < last || d > test.maxDelay {
				t.Errorf("calculateDelayBinary(%d, %s, %s) -> got %d", i, test.baseDelay, test.maxDelay, d)
			}

			last = d
		}
	}
}

func TestCalculateDelay(t *testing.T) {
	maxInt := int((^uint(0)) >> 1)

	tests := []struct {
		min       int
		max       int
		baseDelay time.Duration
		maxDelay  time.Duration
	}{
		{0, 100, time.Millisecond, time.Minute},
		{maxInt - 10, maxInt + 10, time.Nanosecond, time.Minute},
	}

	for _, test := range tests {
		last := time.Duration(0)

		for i := test.min; i < test.max; i++ {
			d := calculateDelay(uint(i), test.baseDelay, test.maxDelay)
			if d < 0 || d < last || d > test.maxDelay {
				t.Errorf("calculateDelay(%d, %s, %s) -> got %d", i, test.baseDelay, test.maxDelay, d)
			}
			last = d
		}
	}
}

func TestMaxRetries(t *testing.T) {
	count := 0

	fn := func() (interface{}, error) {
		count++
		return nil, fmt.Errorf("foo")
	}

	r := Retry(fn, MaxRetries(2))
	r()

	if count != 3 {
		t.Errorf("Want %d, got %d", 3, count)
	}
}

func TestBinaryRaise(t *testing.T) {
	f := func(x uint) bool {
		y := binaryRaise(x)

		if y < 0 {
			return false
		}

		if x <= 62 {
			return y == int64(math.Pow(2, float64(x)))
		} else {
			return y == int64(math.Pow(2, 62))
		}
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

var benchmarkBinaryRaiseResult int64

func BenchmarkBinaryRaise(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		benchmarkBinaryRaiseResult = binaryRaise(uint(n))
	}
}

var benchmarkMathPowResult int64

func BenchmarkMathPowRaise(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		benchmarkBinaryRaiseResult = int64(math.Pow(2, float64(n)))
	}
}

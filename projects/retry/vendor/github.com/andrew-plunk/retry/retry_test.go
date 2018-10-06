package retry

import (
	"errors"
	"testing"

	"github.com/andrew-plunk/retry/strategy"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetryRetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if errorUntilAttemptNumber == attempt {
			return nil
		}

		return errors.New("erroring")
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}

	if errorUntilAttemptNumber != attemptsMade {
		t.Errorf(
			"expected %d attempts to be made, but %d were made instead",
			errorUntilAttemptNumber,
			attemptsMade,
		)
	}
}

func TestShouldAttempt(t *testing.T) {
	shouldAttempt := shouldAttempt(1)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithStrategy(t *testing.T) {
	const attemptNumberShouldReturnFalse = 7

	strategy := func(attempt uint) bool {
		return (attemptNumberShouldReturnFalse != attempt)
	}

	should := shouldAttempt(1, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1+attemptNumberShouldReturnFalse, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(attemptNumberShouldReturnFalse, strategy)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithMultipleStrategies(t *testing.T) {
	trueStrategy := func(attempt uint) bool {
		return true
	}

	falseStrategy := func(attempt uint) bool {
		return false
	}

	should := shouldAttempt(1, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}

type Err struct {
	s []strategy.Strategy
}

func (e Err) Error() string {
	return "Error"
}

func (e Err) Strategies() []strategy.Strategy {
	return e.s
}

func TestRetryable(t *testing.T) {
	var try uint

	retryable := &Err{
		[]strategy.Strategy{
			strategy.Limit(3),
			strategy.Limit(2),
		},
	}
	action := func(t uint) error {
		try = t
		return retryable
	}

	fallback := strategy.Limit(4)
	err := Retry(action, fallback)
	if r, ok := err.(Retryable); !(ok && r == err) {
		t.Errorf("Error %v != %v", err, retryable)
	}

	if try != 2 {
		t.Error("Should have tried 2 times!")
	}

}

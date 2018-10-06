// Package retry provides a simple, stateless, functional mechanism to perform
// actions repetitively until successful.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package retry

import "github.com/andrew-plunk/retry/strategy"

// Action defines a callable function that package retry can handle.
type Action func(attempt uint) error

// Retryable defines an error which specifyies it's retry strategies.
type Retryable interface {
	error
	Strategies() []strategy.Strategy
}

// Retry takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Retry(action Action, strategies ...strategy.Strategy) error {
	return retry(0, action, nil, strategies)
}

func retry(attempt uint, action Action, err error, strategies []strategy.Strategy) error {
	if shouldAttempt(attempt, err, strategies...) {
		if err = action(attempt); err == nil {
			return err
		}

		if attempt == 0 {
			if r, ok := err.(Retryable); ok {
				strategies = r.Strategies()
			}
		}

		attempt++
		return retry(attempt, action, err, strategies)
	}

	return err
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(attempt uint, err error, strategies ...strategy.Strategy) bool {
	shouldAttempt := true

	for i := 0; shouldAttempt && i < len(strategies); i++ {
		shouldAttempt = shouldAttempt && strategies[i](attempt, err)
	}

	return shouldAttempt
}

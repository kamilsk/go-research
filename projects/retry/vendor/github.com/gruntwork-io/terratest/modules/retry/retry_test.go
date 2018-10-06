package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoWithRetry(t *testing.T) {
	t.Parallel()

	expectedOutput := "expected"
	expectedError := fmt.Errorf("expected error")

	actionAlwaysReturnsExpected := func() (string, error) { return expectedOutput, nil }
	actionAlwaysReturnsError := func() (string, error) { return expectedOutput, expectedError }

	createActionThatReturnsExpectedAfterFiveRetries := func() func() (string, error) {
		count := 0
		return func() (string, error) {
			count++
			if count > 5 {
				return expectedOutput, nil
			}
			return expectedOutput, expectedError
		}
	}

	testCases := []struct {
		description   string
		maxRetries    int
		expectedError error
		action        func() (string, error)
	}{
		{"Return value on first try", 10, nil, actionAlwaysReturnsExpected},
		{"Return error on all retries", 10, MaxRetriesExceeded{Description: "Return error on all retries", MaxRetries: 10}, actionAlwaysReturnsError},
		{"Return value after 5 retries", 10, nil, createActionThatReturnsExpectedAfterFiveRetries()},
		{"Return value after 5 retries, but only do 4 retries", 4, MaxRetriesExceeded{Description: "Return value after 5 retries, but only do 4 retries", MaxRetries: 4}, createActionThatReturnsExpectedAfterFiveRetries()},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := DoWithRetryE(t, testCase.description, testCase.maxRetries, 1*time.Millisecond, testCase.action)
			assert.Equal(t, expectedOutput, actualOutput)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

func TestDoWithTimeout(t *testing.T) {
	t.Parallel()

	expectedOutput := "expected"
	expectedError := fmt.Errorf("expected error")

	actionReturnsValueImmediately := func() (string, error) { return expectedOutput, nil }
	actionReturnsErrorImmediately := func() (string, error) { return "", expectedError }

	createActionThatReturnsValueAfterDelay := func(delay time.Duration) func() (string, error) {
		return func() (string, error) {
			time.Sleep(delay)
			return expectedOutput, nil
		}
	}

	createActionThatReturnsErrorAfterDelay := func(delay time.Duration) func() (string, error) {
		return func() (string, error) {
			time.Sleep(delay)
			return "", expectedError
		}
	}

	testCases := []struct {
		description   string
		timeout       time.Duration
		expectedError error
		action        func() (string, error)
	}{
		{"Returns value immediately", 5 * time.Second, nil, actionReturnsValueImmediately},
		{"Returns error immediately", 5 * time.Second, expectedError, actionReturnsErrorImmediately},
		{"Returns value after 2 seconds", 5 * time.Second, nil, createActionThatReturnsValueAfterDelay(2 * time.Second)},
		{"Returns error after 2 seconds", 5 * time.Second, expectedError, createActionThatReturnsErrorAfterDelay(2 * time.Second)},
		{"Returns value after timeout exceeded", 5 * time.Second, TimeoutExceeded{Description: "Returns value after timeout exceeded", Timeout: 5 * time.Second}, createActionThatReturnsValueAfterDelay(10 * time.Second)},
		{"Returns error after timeout exceeded", 5 * time.Second, TimeoutExceeded{Description: "Returns error after timeout exceeded", Timeout: 5 * time.Second}, createActionThatReturnsErrorAfterDelay(10 * time.Second)},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := DoWithTimeoutE(t, testCase.description, testCase.timeout, testCase.action)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

func TestDoInBackgroundUntilStopped(t *testing.T) {
	t.Parallel()

	sleepBetweenRetries := 5 * time.Second
	counter := 0

	stop := DoInBackgroundUntilStopped(t, t.Name(), sleepBetweenRetries, func() {
		counter++
	})

	time.Sleep(sleepBetweenRetries * 3)
	stop.Done()

	assert.Equal(t, 3, counter)

	time.Sleep(sleepBetweenRetries * 3)
	assert.Equal(t, 3, counter)
}

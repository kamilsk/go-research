package retry_test

import (
	"fmt"
	"log"
	"time"

	"github.com/bernos/go-retry"
)

func Example() {
	count := 0

	fn := func() (interface{}, error) {
		count++

		if count < 3 {
			fmt.Println(count)
			return nil, fmt.Errorf("gimme at least 2")
		}
		return "Thats more like it!", nil
	}

	r := retry.Retry(
		fn,
		retry.MaxRetries(5),
		retry.BaseDelay(time.Nanosecond))

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output:
	// 1
	// 2
	// Thats more like it!, <nil>
}

func ExampleRetry() {
	count := 0

	fn := func() (interface{}, error) {
		if count < 2 {
			count++
			fmt.Println(count)
			return nil, fmt.Errorf("gimme at least 2")
		}
		return "Thats more like it!", nil
	}

	r := retry.Retry(fn)

	v, err := r()

	fmt.Printf("%s, %v", v, err)

	// Output:
	// 1
	// 2
	// Thats more like it!, <nil>
}

func ExampleBackoff() {
	fn := func() (interface{}, error) {
		return "foo", nil
	}

	backoff := func(iteration uint, baseDelay, maxDelay time.Duration) time.Duration {
		d := time.Duration(iteration) * baseDelay

		if d < maxDelay {
			return d
		}

		return maxDelay
	}

	r := retry.Retry(fn, retry.Backoff(backoff))

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output: foo, <nil>
}

func ExampleBaseDelay() {
	fn := func() (interface{}, error) {
		return "foo", nil
	}

	r := retry.Retry(fn, retry.BaseDelay(time.Millisecond))

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output: foo, <nil>
}

func ExampleExponentialBackoff() {
	fn := func() (interface{}, error) {
		return "foo", nil
	}

	r := retry.Retry(fn, retry.BinaryExponentialBackoff())

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output: foo, <nil>
}

func ExampleForever() {
	fn := func() (interface{}, error) {
		return "foo", nil
	}

	r := retry.Retry(fn, retry.Forever())

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output: foo, <nil>
}

func ExampleLog() {
	fn := func() (interface{}, error) {
		return "foo", nil
	}

	r := retry.Retry(fn, retry.Log(log.Printf))

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output: foo, <nil>
}

func ExampleMaxDelay() {
	fn := func() (interface{}, error) {
		return "foo", nil
	}

	r := retry.Retry(fn, retry.MaxDelay(time.Minute))

	value, err := r()

	fmt.Printf("%s, %v", value, err)

	// Output: foo, <nil>
}

func ExampleMaxRetries() {
	count := 0

	fn := func() (interface{}, error) {
		count++
		fmt.Println(count)
		return nil, fmt.Errorf("err")
	}

	r := retry.Retry(fn, retry.MaxRetries(3))

	value, err := r()

	fmt.Printf("%v, %v", value, err)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// <nil>, Retrier exceeded max retry count of 3. Cause: err
}

func ExampleWhile() {
	count := 0

	fn := func() (interface{}, error) {
		count++
		fmt.Println(count)
		return nil, fmt.Errorf("err")
	}

	shouldRetry := func(err error) bool {
		return count < 3
	}

	r := retry.Retry(fn, retry.While(shouldRetry))

	value, err := r()

	fmt.Printf("%v, %v", value, err)

	// Output:
	// 1
	// 2
	// 3
	// <nil>, Retrier aborted due to user supplied ShouldRetry func. Cause: err
}

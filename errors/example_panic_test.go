package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/oxequa/grace"
)

func Example_usage() {
	origin := errors.New("origin")

	var err error

	func() {
		defer grace.Recover(&err).Error() // without stack trace
		panic(origin)
	}()

	if reflect.DeepEqual(origin, err) {
		panic("unexpected equality")
	}
	if !grace.Equal(origin, err) {
		panic("equality is expected, but it is true only for `grace.Recover(&err).Error()`")
	}

	fmt.Println("it is not possible to obtain original error")
	// Output: it is not possible to obtain original error
}

func Benchmark_Recover(b *testing.B) {
	var err error
	text := "error"
	b.Run("built-in recover", func(b *testing.B) {
		b.ReportAllocs()
		testCase := func() {
			defer func() {
				if r := recover(); r != nil {
					switch e := (r).(type) {
					case error:
						err = e
					}
				}
			}()
			panic(errors.New(text))
		}
		for i := 0; i < b.N; i++ {
			testCase()
			err = nil
		}
	})
	b.Run("github.com/oxequa/grace", func(b *testing.B) {
		b.ReportAllocs()
		testCase := func() {
			defer grace.Recover(&err).Error()
			panic(errors.New(text))
		}
		for i := 0; i < b.N; i++ {
			testCase()
			err = nil
		}
	})
}

package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/oxequa/grace"
)

func Example_usage() {
	var err error
	origin := errors.New("origin")

	func() {
		defer grace.Recover(&err).Error()
		panic(origin)
	}()
	_, _ = fmt.Println(err.Error())

	if reflect.DeepEqual(origin, err) {
		panic("unexpected equality")
	}
	if !grace.Equal(origin, err) {
		panic("equality is expected, but it is true only for `grace.Recover(&err).Error()`")
	}
	_, _ = fmt.Println("it is not possible to obtain original error")

	// Output:
	// origin
	// it is not possible to obtain original error
}

func Benchmark_Recover(b *testing.B) {
	text := "error"
	b.Run("built-in recover", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		test := func() {
			defer func(err *error) {
				if r := recover(); r != nil {
					switch e := (r).(type) {
					case error:
						*err = e
					}
				}
			}(&err)
			panic(errors.New(text))
		}
		for i := 0; i < b.N; i++ {
			test()
		}
	})
	b.Run("github.com/oxequa/grace", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		test := func() {
			defer grace.Recover(&err).Error()
			panic(errors.New(text))
		}
		for i := 0; i < b.N; i++ {
			test()
		}
	})
}

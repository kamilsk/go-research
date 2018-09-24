package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	deep "github.com/pkg/errors"

	"github.com/goph/emperror"
	"github.com/oxequa/grace"
)

func Example_emperrorUsage() {
	var err error
	origin := errors.New("origin")

	func() {
		defer emperror.HandleRecover(emperror.HandlerFunc(func(handled error) { err = handled }))
		panic(origin)
	}()
	_, _ = fmt.Println(err.Error())

	if reflect.DeepEqual(origin, err) {
		panic("unexpected equality")
	}
	if !reflect.DeepEqual(origin, deep.Cause(err)) {
		panic("equality is expected")
	}
	_, _ = fmt.Println("it is possible to obtain original error by `errors.Cause()`")

	// Output:
	// origin
	// it is possible to obtain original error by `errors.Cause()`
}

func Example_graceUsage() {
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
	origin := errors.New("origin")
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
			panic(origin)
		}
		for i := 0; i < b.N; i++ {
			test()
		}
	})
	b.Run("github.com/goph/emperror", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		handler := func(err *error) emperror.HandlerFunc {
			return func(handled error) {
				*err = handled
			}
		}(&err)
		test := func() {
			defer emperror.HandleRecover(handler)
			panic(origin)
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
			panic(origin)
		}
		for i := 0; i < b.N; i++ {
			test()
		}
	})
}

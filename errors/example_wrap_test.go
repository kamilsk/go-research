package errors

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"testing"

	juju "github.com/juju/errors"
	pkg "github.com/pkg/errors"
)

var showStackTrace = flag.Bool("sst", false, "print error stack trace")

func Example_difference() {
	origin := func() error { return errors.New("origin error") }

	var err error

	message := "errors replacer"
	err = juju.New(message) // has stack trace
	err = pkg.New(message)  // has stack trace

	message = "fmt.Errorf replacer"
	err = juju.Errorf(message) // has stack trace
	err = pkg.Errorf(message)  // has stack trace

	// add only context
	message = "with context"
	err = pkg.WithMessage(origin(), message)

	// add only stack trace
	err = pkg.WithStack(origin())

	// add context and stack trace
	message = "with context and stack trace"
	err = juju.Annotate(origin(), message)
	err = pkg.Wrap(origin(), message)

	func() {
		raw := [1024]byte{}
		buf := bytes.NewBuffer(raw[:0])

		err = func() error { return juju.Annotate(origin(), message) }()
		fmt.Fprintf(buf, "github.com/juju/errors:\n%+v\n", err)
		fmt.Fprintln(buf)
		err = func() error { return pkg.Wrap(origin(), message) }()
		fmt.Fprintf(buf, "github.com/pkg/errors:\n%+v\n", err)

		if *showStackTrace {
			fmt.Println(buf.String())
		}
	}()

	message = "github.com/pkg/errors is my preferred choice"
	fmt.Println(message)
	// Output: github.com/pkg/errors is my preferred choice
}

func Benchmark_New(b *testing.B) {
	message := "error"
	b.Run("github.com/juju/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			juju.New(message)
		}
	})
	b.Run("github.com/pkg/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pkg.New(message)
		}
	})
}

func Benchmark_Wrap(b *testing.B) {
	message := "with context and stack trace"
	origin := errors.New("error")
	b.Run("github.com/juju/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			juju.Annotate(origin, message)
		}
	})
	b.Run("github.com/pkg/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pkg.Wrap(origin, message)
		}
	})
}

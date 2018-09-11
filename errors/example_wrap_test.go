package errors

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	juju "github.com/juju/errors"
	pkg "github.com/pkg/errors"
)

func Example_difference() {
	var err error
	origin := func() error { return errors.New("origin error") }
	buf := bytes.NewBuffer(nil)

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
		err = func() error { return juju.Annotate(origin(), message) }()
		_, _ = fmt.Fprintf(buf, "github.com/juju/errors:\n%+v\n", err)
		_, _ = fmt.Fprintln(buf)
		err = func() error { return pkg.Wrap(origin(), message) }()
		_, _ = fmt.Fprintf(buf, "github.com/pkg/errors:\n%+v\n", err)
	}()

	// sanitize the result https://github.com/golang/go/issues/18831
	result := buf.String()
	result = strings.Replace(result, filepath.Join(build.Default.GOPATH, "src")+string(filepath.Separator), "", -1)
	result = strings.Replace(result, filepath.Join(runtime.GOROOT(), "src")+string(filepath.Separator), "", -1)
	_, _ = fmt.Println(result)
	// Output:
	// github.com/juju/errors:
	// origin error
	// github.com/kamilsk/go-research/errors/example_wrap_test.go:43: with context and stack trace
	//
	// github.com/pkg/errors:
	// origin error
	// with context and stack trace
	// github.com/kamilsk/go-research/errors.Example_difference.func2.2
	//	github.com/kamilsk/go-research/errors/example_wrap_test.go:46
	// github.com/kamilsk/go-research/errors.Example_difference.func2
	//	github.com/kamilsk/go-research/errors/example_wrap_test.go:46
	// github.com/kamilsk/go-research/errors.Example_difference
	//	github.com/kamilsk/go-research/errors/example_wrap_test.go:48
	// testing.runExample
	//	testing/example.go:121
	// testing.runExamples
	//	testing/example.go:45
	// testing.(*M).Run
	//	testing/testing.go:1035
	// main.main
	//	_testmain.go:50
	// runtime.main
	//	runtime/proc.go:201
	// runtime.goexit
	//	runtime/asm_amd64.s:1333
}

func Benchmark_New(b *testing.B) {
	message := "error"
	b.Run("github.com/juju/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = juju.New(message)
		}
	})
	b.Run("github.com/pkg/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = pkg.New(message)
		}
	})
}

func Benchmark_Wrap(b *testing.B) {
	origin, message := errors.New("error"), "with context and stack trace"
	b.Run("github.com/juju/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = juju.Annotate(origin, message)
		}
	})
	b.Run("github.com/pkg/errors", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = pkg.Wrap(origin, message)
		}
	})
}

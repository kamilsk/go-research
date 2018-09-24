package logging

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stdLogger interface {
	SetOutput(io.Writer)
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

var (
	_ stdLogger = log.New(nil, "", 0)
	_ stdLogger = logrus.New()

	sanitizer = regexp.MustCompile(`\d{4}[-/]\d{2}[-/]\d{2}[T\s]\d{2}:\d{2}:\d{2}(?:(?:\.\d{3})?\+\d{2}:?\d{2})?`)
)

func TestRace(t *testing.T) {
	tests := []struct {
		name   string
		writer func(starter <-chan struct{}) func(waiter *sync.WaitGroup, log string)
	}{
		{"github.com/sirupsen/logrus", func(starter <-chan struct{}) func(waiter *sync.WaitGroup, log string) {
			logger := logrus.New()
			logger.SetFormatter(&logrus.TextFormatter{})
			logger.SetLevel(logrus.ErrorLevel)
			logger.SetOutput(ioutil.Discard)
			return func(waiter *sync.WaitGroup, log string) {
				<-starter
				logger.Error(log)
				waiter.Done()
			}
		}},
		{"go.uber.org/zap", func(starter <-chan struct{}) func(waiter *sync.WaitGroup, log string) {
			logger := zap.New(
				zapcore.NewCore(
					zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
						MessageKey:     "msg",
						LevelKey:       "level",
						TimeKey:        "time",
						NameKey:        "logger",
						CallerKey:      "caller",
						StacktraceKey:  "stacktrace",
						LineEnding:     zapcore.DefaultLineEnding,
						EncodeLevel:    zapcore.LowercaseLevelEncoder,
						EncodeTime:     zapcore.ISO8601TimeEncoder,
						EncodeDuration: zapcore.SecondsDurationEncoder,
						EncodeCaller:   zapcore.ShortCallerEncoder,
					}),
					zapcore.AddSync(ioutil.Discard),
					zapcore.ErrorLevel,
				),
			)
			return func(waiter *sync.WaitGroup, log string) {
				<-starter
				logger.Error(log)
				waiter.Done()
			}
		}},
		{"github.com/rs/zerolog", func(starter <-chan struct{}) func(waiter *sync.WaitGroup, log string) {
			logger := zerolog.New(zerolog.ConsoleWriter{Out: ioutil.Discard, NoColor: true}).With().Timestamp().Logger()
			logger.Level(zerolog.ErrorLevel)
			return func(waiter *sync.WaitGroup, log string) {
				<-starter
				logger.Error().Msg(log)
				waiter.Done()
			}
		}},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			waiter := &sync.WaitGroup{}
			starter := make(chan struct{})
			writer := tc.writer(starter)
			for i := 0; i < 100; i++ {
				waiter.Add(1)
				go writer(waiter, fmt.Sprintf("msg#%03d", i))
			}
			close(starter)
			waiter.Wait()
		})
	}
}

func Example_logrusUsage() {
	buf := bytes.NewBuffer(nil)

	// instantiation and configuration
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetOutput(buf)

	// integration with built-in logger
	w := logger.Writer()
	l := log.New(w, "", 0)
	l.Println("built-in logger uses github.com/sirupsen/logrus as writer")
	_ = w.Close()

	// nested logger
	entry := logger.WithFields(logrus.Fields{"default": "value"})
	entry.WithField("key", "value").Info("ignored")

	// usage
	logger.WithField("logger", "github.com/sirupsen/logrus").Error("something happen")

	// sanitize the result https://github.com/golang/go/issues/18831
	result := buf.String()
	result = sanitizer.ReplaceAllString(result, time.Time{}.Format(time.RFC3339))
	_, _ = fmt.Println(result)
	// Output:
	// time="0001-01-01T00:00:00Z" level=error msg="something happen" logger=github.com/sirupsen/logrus
}

func Example_zapUsage() {
	buf := bytes.NewBuffer(nil)

	// instantiation and configuration
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				MessageKey:     "msg",
				LevelKey:       "level",
				TimeKey:        "time",
				NameKey:        "logger",
				CallerKey:      "caller",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(buf),
			zapcore.ErrorLevel,
		),
	)

	// nested logger
	entry := logger.With(zap.String("default", "value"))
	entry.With(zap.String("key", "value")).Info("ignored")

	// usage
	logger.With(zap.String("logger", "go.uber.org/zap")).Error("something happen")
	_ = logger.Sync()

	// sanitize the result https://github.com/golang/go/issues/18831
	result := buf.String()
	result = sanitizer.ReplaceAllString(result, time.Time{}.Format(time.RFC3339))
	_, _ = fmt.Println(result)
	// Output:
	// 0001-01-01T00:00:00Z	error	something happen	{"logger": "go.uber.org/zap"}
}

func Example_zerologUsage() {
	buf := bytes.NewBuffer(nil)

	// instantiation and configuration
	logger := zerolog.New(zerolog.ConsoleWriter{Out: buf, NoColor: true}).With().Timestamp().Logger()
	logger.Level(zerolog.ErrorLevel)

	// integration with built-in logger
	l := log.New(logger, "", 0)
	l.Println("built-in logger uses github.com/rs/zerolog as writer")

	// nested logger
	entry := logger.With().Str("default", "value").Logger().Level(zerolog.ErrorLevel)
	entry.Info().Str("key", "value").Msg("ignored")

	// usage
	logger.Error().Str("logger", "github.com/rs/zerolog").Msg("something happen")

	// sanitize the result https://github.com/golang/go/issues/18831
	result := buf.String()
	result = sanitizer.ReplaceAllString(result, time.Time{}.Format(time.RFC3339))
	_, _ = fmt.Println(result)
	// Output:
	// 0001-01-01T00:00:00Z |????| built-in logger uses github.com/rs/zerolog as writer
	// 0001-01-01T00:00:00Z |ERROR| something happen logger=github.com/rs/zerolog
}

func Benchmark_Usage(b *testing.B) {
	b.Run("github.com/sirupsen/logrus", func(b *testing.B) {
		b.ReportAllocs()
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{})
		logger.SetLevel(logrus.InfoLevel)
		logger.SetOutput(ioutil.Discard)
		for i := 0; i < b.N; i++ {
			logger.WithFields(logrus.Fields{
				"logger":  "github.com/sirupsen/logrus",
				"error":   "context",
				"package": "name",
				"file":    "path",
				"line":    "number",
			}).Error("something happen")
		}
	})
	b.Run("go.uber.org/zap", func(b *testing.B) {
		b.ReportAllocs()
		logger := zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
					MessageKey:     "msg",
					LevelKey:       "level",
					TimeKey:        "time",
					NameKey:        "logger",
					CallerKey:      "caller",
					StacktraceKey:  "stacktrace",
					LineEnding:     zapcore.DefaultLineEnding,
					EncodeLevel:    zapcore.LowercaseLevelEncoder,
					EncodeTime:     zapcore.ISO8601TimeEncoder,
					EncodeDuration: zapcore.SecondsDurationEncoder,
					EncodeCaller:   zapcore.ShortCallerEncoder,
				}),
				zapcore.AddSync(ioutil.Discard),
				zapcore.InfoLevel,
			),
		)
		for i := 0; i < b.N; i++ {
			logger.With(
				zap.String("logger", "go.uber.org/zap"),
				zap.String("error", "context"),
				zap.String("package", "name"),
				zap.String("file", "path"),
				zap.String("line", "number"),
			).Error("something happen")
		}
	})
	b.Run("github.com/rs/zerolog", func(b *testing.B) {
		b.ReportAllocs()
		logger := zerolog.New(zerolog.ConsoleWriter{Out: ioutil.Discard}).With().Timestamp().Logger()
		logger.Level(zerolog.InfoLevel)
		for i := 0; i < b.N; i++ {
			logger.Error().Fields(map[string]interface{}{
				"logger":  "github.com/rs/zerolog",
				"error":   "context",
				"package": "name",
				"file":    "path",
				"line":    "number",
			}).Msg("something happen")
		}
	})
}

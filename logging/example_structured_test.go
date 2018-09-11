package logging

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
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
	logger := zerolog.New(zerolog.ConsoleWriter{Out: buf}).With().Timestamp().Logger()
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

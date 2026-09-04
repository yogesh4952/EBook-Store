package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const RequestIDKey = "request_id"

type ctxKey struct{}

// Init configures the application-wide zerolog logger.
//
// Behavior is controlled by env vars:
//   - ENV=production   -> JSON logs to stdout (with timestamp)
//   - otherwise        -> pretty, colored console logs to stderr
//   - LOG_LEVEL        -> trace | debug | info | warn | error | fatal
//                        (default: info, so you must set LOG_LEVEL=debug
//                        to see SQL/gorm trace output)
//
// Call this once, as early as possible (before any logging happens).
func Init() {
	var w io.Writer = os.Stdout
	if os.Getenv("ENV") != "production" {
		w = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	}

	level := zerolog.InfoLevel
	if lvl, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil && lvl != zerolog.NoLevel {
		level = lvl
	}

	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(w).Level(level).With().Timestamp().Logger()
}

// Level returns the currently active minimum level.
func Level() string {
	return log.Logger.GetLevel().String()
}

// SetLevel changes the minimum level at runtime, e.g.
//
//	logger.SetLevel("debug")
func SetLevel(level string) {
	if lvl, err := zerolog.ParseLevel(level); err == nil && lvl != zerolog.NoLevel {
		log.Logger = log.Logger.Level(lvl)
	}
}

// Success logs at info level with a "+" marker (used for startup steps).
func Success(format string, a ...any) {
	log.Info().Msgf("[+] "+format, a...)
}

func Info(format string, a ...any) {
	log.Info().Msgf(format, a...)
}

func Warn(format string, a ...any) {
	log.Warn().Msgf(format, a...)
}

func Error(format string, a ...any) {
	log.Error().Msgf(format, a...)
}

func Fatal(format string, a ...any) {
	log.Fatal().Msgf(format, a...)
}

// NewRequestID returns a short random hex id used for request correlation.
func NewRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// WithRequestID attaches the request id to a context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// RequestID extracts the request id from a context ("" if absent).
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(ctxKey{}).(string); ok {
		return id
	}
	return ""
}

// Ctx returns a logger bound to the context. If the context carries a
// request id, it is added to every line, e.g.
//
//	logger.Ctx(ctx).Error().Msg("payment failed")
func Ctx(ctx context.Context) *zerolog.Logger {
	l := log.Logger.With().Logger()
	if rid := RequestID(ctx); rid != "" {
		l = l.With().Str(RequestIDKey, rid).Logger()
	}
	return &l
}
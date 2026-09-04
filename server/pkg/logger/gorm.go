package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gormlogger "gorm.io/gorm/logger"
)

// NewGormLogger wires GORM's output into zerolog.
//
//   - slow    - queries slower than this are logged at WARN ("slow SQL")
//   - level   - the minimum log level GORM is allowed to emit
//     (gormlogger.Warn is the sensible default; use
//     gormlogger.Info + LOG_LEVEL=debug to see every query)
func NewGormLogger(level gormlogger.LogLevel, slow time.Duration) gormlogger.Interface {
	return gormlogger.New(gormWriter{}, gormlogger.Config{
		SlowThreshold:             slow,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
}

type gormWriter struct{}

func (gormWriter) Printf(format string, v ...any) {
	msg := strings.TrimSpace(fmt.Sprintf(format, v...))
	if msg == "" {
		return
	}

	var evt *zerolog.Event
	switch {
	case strings.Contains(msg, "[error]"):
		evt = log.Error()
	case strings.Contains(msg, "[warn]"), strings.Contains(msg, "SLOW SQL"):
		evt = log.Warn()
	default:
		evt = log.Debug()
	}
	evt.Str("component", "gorm").Msg(msg)
}
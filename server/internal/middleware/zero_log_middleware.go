package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

// ZerologMiddleware logs one line per HTTP request.
//
// A request id is read from the X-Request-ID header (or generated) and
// attached both to the log line and to the request context, so downstream
// handlers/services can log with the same id via logger.Ctx(ctx).
//
// Log level follows the response status: 5xx -> error, 4xx -> warn,
// everything else -> info.
func ZerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = logger.NewRequestID()
		}
		c.Set(logger.RequestIDKey, rid)
		c.Request = c.Request.WithContext(logger.WithRequestID(c.Request.Context(), rid))

		c.Next()

		status := c.Writer.Status()
		evt := log.Info()
		switch {
		case status >= 500:
			evt = log.Error()
		case status >= 400:
			evt = log.Warn()
		}

		evt.
			Str(logger.RequestIDKey, rid).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("http request")
	}
}
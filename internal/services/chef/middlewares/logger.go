package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func StructuredLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		logData := slog.Group(
			"data",
			slog.String("client_id", param.ClientIP),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("method", param.Method),
			slog.Int("status_code", param.StatusCode),
			slog.Int("body_size", param.BodySize),
			slog.String("path", param.Path),
			slog.String("latency", param.Latency.String()),
		)

		m := "request"
		if len(param.ErrorMessage) > 0 {
			m = param.ErrorMessage
		}

		if c.Writer.Status() >= 500 {
			logger.Error(m, logData)
		} else {
			logger.Debug(m, logData)
		}
	}
}

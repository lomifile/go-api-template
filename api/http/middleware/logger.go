package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lomifile/api/pkg/logger"
	"go.uber.org/zap"
)

func LoggerMiddleware(l *logger.Logger) fiber.Handler {
	httpLog := l.Named("http")

	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		latency := time.Since(start)

		reqID := c.GetRespHeader(fiber.HeaderXRequestID)
		if reqID == "" {
			reqID = c.Get(fiber.HeaderXRequestID)
		}

		status := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("request_id", reqID),
			zap.String("http_method", c.Method()),
			zap.String("http_path", c.Path()),
			zap.Int("status", status),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.String("client_ip", c.IP()),
			zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
		}

		httpLog.Info("http_request", fields...)

		return err
	}
}

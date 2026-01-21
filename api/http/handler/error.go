package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lomifile/api/pkg/logger"
	"github.com/lomifile/api/pkg/utils"
	"go.uber.org/zap"
)

type ErrorResponder struct {
	l *logger.Logger
}

func NewErrorResponder(l *logger.Logger) *ErrorResponder {
	return &ErrorResponder{l: l}
}

func (r *ErrorResponder) Error(
	c *fiber.Ctx,
	status int,
	message string,
	logKey string,
	fields ...zap.Field,
) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	now := time.Now().String()

	if logKey != "" {
		r.l.Error(
			logKey,
			append([]zap.Field{
				zap.String("request_id", requestID),
				zap.Int("status", status),
			}, fields...)...,
		)
	}

	return c.Status(status).JSON(utils.ErrorResponseMap{
		RequestID: requestID,
		Status:    status,
		Error:     message,
		TS:        now,
	})
}

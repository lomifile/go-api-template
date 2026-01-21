// Package router handles route registration
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lomifile/api/config"
	"github.com/lomifile/api/internal/adapter"
	"github.com/lomifile/api/pkg/logger"
)

func NewRouter(
	app *fiber.App,
	db *adapter.PostgresAdapter,
	l *logger.Logger,
	c *config.Config,
) {
}

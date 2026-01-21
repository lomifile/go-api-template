// Package app handles application bootstrap and lifecycle
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/lomifile/api/api/http/middleware"
	"github.com/lomifile/api/api/http/router"
	"github.com/lomifile/api/config"
	"github.com/lomifile/api/internal/adapter"
	"github.com/lomifile/api/internal/server"
	"github.com/lomifile/api/pkg/logger"
	"github.com/lomifile/api/pkg/postgres"
	"go.uber.org/zap"
)

func Start(c *config.Config) {
	l := logger.New(logger.Config{
		Env:   c.Environment,
		Debug: true,
	})
	defer func() {
		if err := l.Sync(); err != nil {
			panic(err)
		}
	}()

	p, err := postgres.New(c.Database.Dsn)
	if err != nil {
		l.Error("Postgres connection error", zap.String("err", err.Error()))
		panic(err)
	}

	defer p.Close()
	l.Info("Postgres connected successfully")

	db, err := adapter.NewPostgresAdapter(p)
	if err != nil {
		l.Error("Postgres adapter error", zap.String("err", err.Error()))
		panic(err)
	}

	s := server.New(server.Port(c.Port))
	s.App.Use(middleware.LoggerMiddleware(l))
	s.App.Use(requestid.New())
	s.App.Use(recover.New())
	s.App.Use(helmet.New())
	s.App.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Authorization,Content-Type",
		AllowCredentials: true,
	}))

	s.App.Use(encryptcookie.New(encryptcookie.Config{
		Key: c.CookieKey,
	}))
	router.NewRouter(s.App, db, l, c)
	s.Start()

	l.Info(fmt.Sprintf("app started on port %s", c.Port))
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-osSignal:
		l.Info("Signal: ", zap.String("", s.String()))

	case err = <-s.Notify():
		l.Error("Server error: ", zap.String("", err.Error()))
	}

	err = s.Shutdown()
	if err != nil {
		l.Error("Shutdown defect: ", zap.String("", err.Error()))
	}
}

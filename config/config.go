// Package config handles application configuration
package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseOptions struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type LimiterOptions struct {
	RPS     float64
	Burst   int
	Enabled bool
}

type Config struct {
	Port        string
	Environment string
	SecretKey   string
	Database    DatabaseOptions
	Limiter     LimiterOptions
	Email       Email
	CookieKey   string
}

type Email struct {
	Host     string
	Username string
	Password string
}

func New() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	return &Config{}
}

func (c *Config) ParseConfig() error {
	flag.StringVar(&c.Port, "port", os.Getenv("PORT"), "API port")
	flag.StringVar(
		&c.Environment,
		"environment",
		os.Getenv("ENVIRONMENT"),
		"Environment (production|development)",
	)
	flag.StringVar(&c.SecretKey, "secret-key", os.Getenv("SECRET_KEY"), "JWT secret key")
	flag.StringVar(&c.CookieKey, "cookie-eky", os.Getenv("COOKIE_KEY"), "Cookie secret key")

	flag.StringVar(&c.Database.Dsn, "db-dsn", os.Getenv("DB_URL"), "PostgreSQL DSN")

	flag.IntVar(
		&c.Database.MaxOpenConns,
		"db-max-open-conns",
		25,
		"PostgreSQL max open connections",
	)
	flag.IntVar(
		&c.Database.MaxIdleConns,
		"db-max-idle-conns",
		25,
		"PostgreSQL max idle connections",
	)
	flag.StringVar(
		&c.Database.MaxIdleTime,
		"db-max-idle-time",
		"15m",
		"PostgreSQL max connection idle time",
	)

	flag.Float64Var(&c.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&c.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&c.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&c.Email.Host, "email-host", os.Getenv("EMAIL_SERVER"), "Email dsn")
	flag.StringVar(
		&c.Email.Username,
		"email-username",
		os.Getenv("EMAIL_USERNAME"),
		"Email username",
	)
	flag.StringVar(
		&c.Email.Password,
		"email-password",
		os.Getenv("EMAIL_PASSWORD"),
		"Email password",
	)

	return nil
}

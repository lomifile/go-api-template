// Package postgres provides interface for postgres pool
package postgres

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxConns     int32         = 4
	_defaultConnAttempts               = 10
	_defaultConnTimeout  time.Duration = time.Second
)

// Postgres main postgres pool struct
type Postgres struct {
	maxConns     int32
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

// Option Provides function to postgres options
type Option func(*Postgres)

// WithMaxConns Adds maxMaxConns
func WithMaxConns(n int32) Option { return func(p *Postgres) { p.maxConns = n } }

// WithConnAttempts Increments connAttempts
func WithConnAttempts(n int) Option { return func(p *Postgres) { p.connAttempts = n } }

// WithConnTimeout adds conn timeout
func WithConnTimeout(d time.Duration) Option {
	return func(p *Postgres) { p.connTimeout = d }
}

// New creates new postgres pool
func New(dsn string, opts ...Option) (*Postgres, error) {
	p := &Postgres{
		maxConns:     _defaultMaxConns,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}
	for _, opt := range opts {
		opt(p)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse dsn: %w", err)
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return nil, fmt.Errorf("postgres: invalid scheme: %s", u.Scheme)
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse pool cfg: %w", err)
	}
	cfg.MaxConns = p.maxConns

	var pool *pgxpool.Pool
	attempts := p.connAttempts

	for {
		ctx, cancel := context.WithTimeout(context.Background(), p.connTimeout)
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		cancel()
		if err == nil {
			ctx, cancel = context.WithTimeout(context.Background(), p.connTimeout)
			err = pool.Ping(ctx)
			cancel()
		}
		if err == nil {
			break
		}
		attempts--
		if attempts <= 0 {
			return nil, fmt.Errorf("postgres: connection attempts exhausted: %w", err)
		}
		log.Printf("postgres: retrying, attempts left=%d, err=%v", attempts, err)
		time.Sleep(p.connTimeout)
	}

	return &Postgres{
		Pool:         pool,
		maxConns:     p.maxConns,
		connAttempts: p.connAttempts,
		connTimeout:  p.connTimeout,
	}, nil
}

// Close closes postgres pool
func (p *Postgres) Close() {
	if p != nil && p.Pool != nil {
		p.Pool.Close()
	}
}

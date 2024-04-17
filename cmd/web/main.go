package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"nishojib/gotrello/internal/server"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/nedpals/supabase-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		debug        string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     int
		enabled bool
	}
	sb struct {
		url string
		key string
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(
		&cfg.db.dsn,
		"db-dsn",
		"",
		"PostgreSQL DSN",
	)
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-time",
		15*time.Minute,
		"PostgreSQL max connection idle time",
	)
	flag.StringVar(&cfg.db.debug, "db-debug", "", "PostgreSQL debug mode (verbose|)")

	flag.IntVar(&cfg.limiter.rps, "limiter-rps", 100, "Rate limiter requests per second")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.sb.url, "sb-url", "", "Supabase URL")
	flag.StringVar(&cfg.sb.key, "sb-key", "", "Supabase Key")

	flag.Parse()

	bunDB, err := InitBunDB(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	sbClient := supabase.CreateClient(cfg.sb.url, cfg.sb.key)

	s := server.New(
		bunDB,
		sbClient,
		server.NewLimiter(cfg.limiter.rps, cfg.limiter.enabled),
		server.WithPort(cfg.port),
	)

	if err := s.Serve(cfg.env); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func InitBunDB(cfg config) (*bun.DB, error) {
	slog.Info(fmt.Sprintf("Connecting to PostgreSQL with dsn: %s", cfg.db.dsn))
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	bunDB := bun.NewDB(db, pgdialect.New())
	if len(cfg.db.debug) > 0 {
		bunDB.AddQueryHook(
			bundebug.NewQueryHook(bundebug.WithVerbose(true)),
		)
	}

	return bunDB, nil
}

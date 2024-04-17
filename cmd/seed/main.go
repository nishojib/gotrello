package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"nishojib/gotrello/internal/data/models"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type config struct {
	db struct {
		dsn          string
		debug        string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

func main() {
	var cfg config

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

	flag.Parse()

	bunDB, err := InitBunDB(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Running fixtures...")

	if err := runFixtures(bunDB); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Finished running fixtures")
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

func runFixtures(bunDB *bun.DB) error {
	return bunDB.RunInTx(
		context.Background(),
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			bunDB.RegisterModel((*models.Project)(nil), (*models.Status)(nil), (*models.Task)(nil))
			fixture := dbfixture.New(bunDB, dbfixture.WithRecreateTables())
			return fixture.Load(context.Background(), os.DirFS("testdata"), "fixtures/data.yml")
		},
	)
}

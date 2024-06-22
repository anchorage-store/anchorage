/*
usrhttpd runs the public user server.

This server registers users with their public keys for distrubution to
other users to send their vaults and items to.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sethvargo/go-envconfig"

	"github.com/anchorage-store/anchorage/clog"
)

// Holds configuration for running the server.
type config struct {
	Port   int    `env:"USRHTTPD_PORT"`
	DBHost string `env:"USRHTTPD_DB_HOST"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	slog.SetDefault(slog.New(clog.NewHandler(os.Stderr, "user", nil)))

	var (
		c config
		_ = flag.Bool("migrate", false, "--migrate=true")
	)
	if err := envconfig.Process(ctx, &c); err != nil {
		slog.Error("error processing env", "err", err)
		os.Exit(1)
	}

	_, err := connectDB(ctx, c)
	if err != nil {
		slog.Error("error connecting to db", "err", err)
		os.Exit(1)
	}
}

func connectDB(ctx context.Context, c config) (*sqlx.DB, error) {
	dbx, err := sqlx.Open("sqlite3", c.DBHost)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %s", err)
	}

	// Try to ping
	if err := dbx.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %s", err)
	}

	return dbx, nil
}

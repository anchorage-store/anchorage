/*
usrhttpd runs the public user server.

This server registers users with their public keys for distrubution to
other users to send their vaults and items to.
*/
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sethvargo/go-envconfig"

	"github.com/anchorage-store/anchorage/clog"
	"github.com/anchorage-store/anchorage/migration"
	"github.com/anchorage-store/anchorage/serv/usrserv"
)

//go:embed migrations/*
var migrations embed.FS

// Holds configuration for running the server.
type config struct {
	Port   int    `env:"USRHTTPD_PORT"`
	DBHost string `env:"USRHTTPD_DB_HOST"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := realMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	slog.SetDefault(slog.New(clog.NewHandler(os.Stderr, "user", nil)))

	var (
		c       config
		migrate = flag.Bool("migrate", false, "--migrate=true")
	)
	if err := envconfig.Process(ctx, &c); err != nil {
		slog.Error("error processing env", "err", err)
		os.Exit(1)
	}
	flag.Parse()
	slog.Info("configuration parsed", "config", c)

	db, err := connectDB(ctx, c)
	if err != nil {
		return err
	}

	// If migrate = true, run migrations instead of the server
	if *migrate {
		applied, err := migration.Migrate(ctx, db, migrations)
		if err != nil {
			return err
		}

		slog.Info("migrated", "applied", applied)
		return nil
	}

	errs := make(chan error)
	s := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.Port),
		Handler:      usrserv.Router(),
		WriteTimeout: time.Second,
		ReadTimeout:  time.Second,
	}
	go func() {
		// Starting this in a goroutine so the main routine can wait
		// for context cancellation
		slog.Info("server listening", "on", c.Port)
		if err := s.ListenAndServe(); err != nil {
			errs <- fmt.Errorf("error serving: %s", err)
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errs:
		return err
	}

	return nil
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

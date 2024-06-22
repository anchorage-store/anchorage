/*
usrhttpd runs the public user server.

This server registers users with their public keys for distrubution to
other users to send their vaults and items to.
*/
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/anchorage-store/anchorage/clog"
	"github.com/sethvargo/go-envconfig"
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
}

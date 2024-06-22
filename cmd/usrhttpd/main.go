/*
usrhttpd runs the public user server.

This server registers users with their public keys for distrubution to
other users to send their vaults and items to.
*/
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

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

	var c config
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
}

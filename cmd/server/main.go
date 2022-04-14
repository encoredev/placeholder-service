package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.encore.dev/placeholder-service/pkg/config"
	"go.encore.dev/placeholder-service/pkg/http"
)

func main() {
	// Start the main context for the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialise our logging library
	log.Logger = zerolog.New(
		zerolog.NewConsoleWriter(),
	).With().Caller().Timestamp().Logger()

	// Listen for OS level signals to shutdown and then cancel our main context
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-done
		log.Warn().Str("signal", s.String()).Msg("received signal to shutdown")
		cancel()
	}()

	// Load the config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	// Start the HTTP server
	err = http.Listen(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred while serving the http server")
		os.Exit(2)
	}

	log.Info().Msg("server shutdown cleanly")
}

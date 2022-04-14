package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.encore.dev/placeholder-service/pkg/config"
)

func Listen(ctx context.Context, cfg *config.Config) error {
	// Setup the router
	var router = mux.NewRouter()
	router.Use(PanicRecovery(), RequestLogger())

	for _, route := range cfg.HealthCheckRoutes {
		log.Info().Str("url", route).Msg("adding health check route")
		router.Methods("GET").Path(route).Handler(http.HandlerFunc(handleHealthRoute))
	}
	router.Methods("GET").PathPrefix("/").Handler(http.HandlerFunc(handleOtherRoutes))

	// Start the server
	log.Info().Int("port", cfg.HttpPort).Msg("starting http server")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: router,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ConnContext: func(ctx context.Context, _ net.Conn) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		log.Warn().Msg("shutting down http server")
		if err := srv.Close(); err != nil {
			log.Err(err).Msg("error shutting down http server")
		}
	}()

	err := srv.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return errors.Wrap(err, "error listening to http server")
	}

	return nil
}

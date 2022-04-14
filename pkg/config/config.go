package config

import (
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort          int      // What port should we be listening on?
	HealthCheckRoutes []string // what URLs should we host health check responses for?
}

func Load() (*Config, error) {
	// Load the .env file if present
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, errors.Wrap(err, "unable to load env")
	}

	// Now configure viper with our default config and bind it to read from the environment
	viper.SetDefault("http_port", 8080)
	viper.SetDefault("health_check_routes", []string{"/healthz", "/__encore/healthz"})
	viper.AutomaticEnv()

	// Read the config file
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/encore/placeholder-service")
	viper.AddConfigPath("$HOME/.encore/placeholder-service")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		// Ignore file not found errors
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "unable to read config.yaml")
		}
	}

	// Read the config
	httpPort := viper.GetInt("http_port")
	if httpPort <= 0 {
		return nil, errors.New("http_port must be greater than 0")
	}

	healthCheckRoutes := viper.GetStringSlice("health_check_routes")
	for i, route := range healthCheckRoutes {
		// Ensure all routes start with a slash
		if !strings.HasPrefix(route, "/") {
			healthCheckRoutes[i] = "/" + route
		}
	}

	log.Info().Int("http_port", httpPort).
		Int("num_health_routes", len(healthCheckRoutes)).
		Msg("loaded service config")

	return &Config{
		HttpPort:          httpPort,
		HealthCheckRoutes: healthCheckRoutes,
	}, nil
}

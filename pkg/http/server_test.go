package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/frankban/quicktest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.encore.dev/placeholder-service/pkg/config"
)

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(
		zerolog.NewConsoleWriter(),
	).With().Caller().Timestamp().Logger()

	os.Exit(m.Run())
}

func Test_Listen(t *testing.T) {
	c := quicktest.New(t)
	c.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a basic config
	cfg := &config.Config{
		HttpPort:          mustFreePort(c),
		HealthCheckRoutes: []string{"/health-route-1", "/another-health-route"},
	}

	// Start Server
	serverShutdown := make(chan error)
	go func() {
		serverShutdown <- Listen(ctx, cfg)
	}()

	// Test we get a 404 on the root route
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", cfg.HttpPort))
	c.Assert(err, quicktest.IsNil)
	c.Assert(resp.StatusCode, quicktest.Equals, http.StatusNotFound)
	_ = resp.Body.Close()

	// Test the first health route
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/health-route-1", cfg.HttpPort))
	c.Assert(err, quicktest.IsNil)
	c.Assert(resp.StatusCode, quicktest.Equals, http.StatusOK)
	_ = resp.Body.Close()

	// Test the same URL but with a suffix on it - we expect a 404
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/health-route-1/", cfg.HttpPort))
	c.Assert(err, quicktest.IsNil)
	c.Assert(resp.StatusCode, quicktest.Equals, http.StatusNotFound)
	_ = resp.Body.Close()

	// Test the other health route
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/another-health-route", cfg.HttpPort))
	c.Assert(err, quicktest.IsNil)
	c.Assert(resp.StatusCode, quicktest.Equals, http.StatusOK)
	_ = resp.Body.Close()

	// Shut the server down
	cancel()
	c.Assert(<-serverShutdown, quicktest.IsNil, quicktest.Commentf("listen returned error"))
}

func mustFreePort(c *quicktest.C) int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		c.Fatalf("unable to get free port %+v", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		c.Fatalf("unable to get free port %+v", err)
	}
	defer func() { _ = l.Close() }()

	return l.Addr().(*net.TCPAddr).Port
}

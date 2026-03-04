// Command unifi-mcp is an MCP server that exposes UniFi network management
// operations as tools for use with AI assistants and MCP clients.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gordcurrie/unifi-mcp/internal/unifi"
	"github.com/gordcurrie/unifi-mcp/tools"
)

// version is set at build time via -ldflags "-X main.version=<tag>".
// Falls back to "dev" when built without ldflags (e.g. go run or go build locally).
var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	var transport string
	var addr string
	flag.StringVar(&transport, "transport", "stdio", "transport to use: stdio or http")
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "listen address for http transport")
	flag.Parse()

	baseURL := os.Getenv("UNIFI_BASE_URL")
	apiKey := os.Getenv("UNIFI_API_KEY")
	siteID := os.Getenv("UNIFI_SITE_ID")
	insecure := os.Getenv("UNIFI_INSECURE") == "true"

	allowDestructive := os.Getenv("UNIFI_ALLOW_DESTRUCTIVE") == "true"

	client, err := unifi.NewClient(baseURL, apiKey, siteID, insecure)
	if err != nil {
		return fmt.Errorf("unifi client: %w", err)
	}

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "unifi-mcp",
		Version: version,
	}, nil)

	tools.RegisterAll(s, client, allowDestructive)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch transport {
	case "stdio":
		if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil && !errors.Is(err, context.Canceled) {
			return fmt.Errorf("stdio transport: %w", err)
		}
	case "http":
		httpServer := &http.Server{
			Addr:              addr,
			Handler:           http.MaxBytesHandler(mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server { return s }, nil), 4<<20),
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1 MiB
		}
		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				slog.Error("http server shutdown", "err", err)
			}
		}()
		slog.Warn("HTTP transport has no authentication — restrict network access to trusted hosts only")
		slog.Info("unifi-mcp listening", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
	default:
		return fmt.Errorf("unknown transport %q (use stdio or http)", transport)
	}
	return nil
}

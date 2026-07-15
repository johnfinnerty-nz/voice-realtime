package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lixuanqun/voice-realtime/internal/config"
	"github.com/lixuanqun/voice-realtime/internal/gateway"
	"github.com/lixuanqun/voice-realtime/internal/providers"

	_ "github.com/lixuanqun/voice-realtime/internal/providers/bailian"
	_ "github.com/lixuanqun/voice-realtime/internal/providers/stepfun"
	_ "github.com/lixuanqun/voice-realtime/internal/providers/volcengine"
	_ "github.com/lixuanqun/voice-realtime/internal/providers/zhipu"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	srv := gateway.New(cfg, slog.Default())
	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("voice-realtime listening", "addr", cfg.Addr, "providers", providers.List())
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	slog.Info("shutdown complete")
}

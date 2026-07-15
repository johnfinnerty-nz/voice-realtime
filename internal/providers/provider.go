package providers

import (
	"context"

	"github.com/lixuanqun/voice-realtime/internal/config"
)

// ConnectConfig is passed when dialing an upstream provider.
type ConnectConfig struct {
	Model  string
	Config *config.Config
}

// UpstreamConn is a bidirectional connection to a cloud voice backend.
type UpstreamConn interface {
	WriteClientEvent(ctx context.Context, raw []byte) error
	ReadServerEvent(ctx context.Context) ([]byte, error)
	Close() error
}

// Provider dials a cloud backend and returns an UpstreamConn.
type Provider interface {
	Name() string
	Dial(ctx context.Context, cfg ConnectConfig) (UpstreamConn, error)
}

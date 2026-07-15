package stepfun

import (
	"context"
	"net/http"

	"github.com/lixuanqun/voice-realtime/internal/providers"
	"github.com/lixuanqun/voice-realtime/internal/providers/proxy"
)

const (
	defaultURL   = "wss://api.stepfun.com/v1/realtime"
	defaultModel = "stepaudio-2.5-realtime"
)

// Provider connects to StepFun Realtime API (OpenAI Realtime compatible).
type Provider struct{}

func New() *Provider { return &Provider{} }

func init() {
	providers.Register(New())
}

func (p *Provider) Name() string { return "stepfun" }

func (p *Provider) Dial(ctx context.Context, cfg providers.ConnectConfig) (providers.UpstreamConn, error) {
	model := cfg.Model
	if model == "" {
		model = defaultModel
	}
	hdr := http.Header{}
	hdr.Set("Authorization", "Bearer "+cfg.Config.StepFunAPIKey)
	return proxy.Dial(ctx, proxy.Options{
		BaseURL: defaultURL,
		Headers: hdr,
		Model:   model,
	})
}

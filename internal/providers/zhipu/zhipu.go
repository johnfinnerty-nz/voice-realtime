package zhipu

import (
	"context"
	"net/http"

	"github.com/lixuanqun/voice-realtime/internal/providers"
	"github.com/lixuanqun/voice-realtime/internal/providers/proxy"
)

const (
	defaultURL   = "wss://open.bigmodel.cn/api/paas/v4/realtime"
	defaultModel = "glm-realtime-flash"
)

// Provider connects to Zhipu GLM-Realtime (OpenAI Realtime compatible).
type Provider struct{}

func New() *Provider { return &Provider{} }

func init() {
	providers.Register(New())
}

func (p *Provider) Name() string { return "zhipu" }

func (p *Provider) Dial(ctx context.Context, cfg providers.ConnectConfig) (providers.UpstreamConn, error) {
	model := cfg.Model
	if model == "" {
		model = defaultModel
	}
	hdr := http.Header{}
	hdr.Set("Authorization", "Bearer "+cfg.Config.ZhipuAPIKey)
	return proxy.Dial(ctx, proxy.Options{
		BaseURL: defaultURL,
		Headers: hdr,
		Model:   model,
	})
}

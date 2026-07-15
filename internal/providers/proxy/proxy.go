package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/coder/websocket"
	"github.com/lixuanqun/voice-realtime/internal/providers"
)

// Options configures a transparent WebSocket proxy to an OpenAI-compatible upstream.
type Options struct {
	BaseURL string
	Headers http.Header
	// ModelQuery adds ?model= if set and not already in URL.
	Model string
}

// Conn proxies JSON events between client and upstream Realtime API.
type Conn struct {
	ws *websocket.Conn
}

// Dial connects to the upstream Realtime WebSocket endpoint.
func Dial(ctx context.Context, opts Options) (*Conn, error) {
	u, err := url.Parse(opts.BaseURL)
	if err != nil {
		return nil, err
	}
	if opts.Model != "" && u.Query().Get("model") == "" {
		q := u.Query()
		q.Set("model", opts.Model)
		u.RawQuery = q.Encode()
	}

	hdr := http.Header{}
	for k, vs := range opts.Headers {
		for _, v := range vs {
			hdr.Add(k, v)
		}
	}

	ws, resp, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPHeader: hdr,
	})
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, fmt.Errorf("upstream dial: %w", err)
	}
	return &Conn{ws: ws}, nil
}

func (c *Conn) WriteClientEvent(ctx context.Context, raw []byte) error {
	return c.ws.Write(ctx, websocket.MessageText, raw)
}

func (c *Conn) ReadServerEvent(ctx context.Context) ([]byte, error) {
	typ, data, err := c.ws.Read(ctx)
	if err != nil {
		return nil, err
	}
	if typ != websocket.MessageText {
		return nil, fmt.Errorf("unexpected upstream message type %d", typ)
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func (c *Conn) Close() error {
	return c.ws.Close(websocket.StatusNormalClosure, "session closed")
}

var _ providers.UpstreamConn = (*Conn)(nil)

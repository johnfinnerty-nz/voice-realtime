package volcengine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/lixuanqun/voice-realtime/internal/providers"
	"github.com/lixuanqun/voice-realtime/internal/realtime"
)

// Provider adapts Volcengine Doubao realtime dialogue to OpenAI Realtime events.
type Provider struct{}

func New() *Provider { return &Provider{} }

func init() {
	providers.Register(New())
}

func (p *Provider) Name() string { return "volcengine" }

func (p *Provider) Dial(ctx context.Context, cfg providers.ConnectConfig) (providers.UpstreamConn, error) {
	connectID := newID()
	sessionID := newID()

	hdr := http.Header{}
	hdr.Set("X-Api-App-ID", cfg.Config.VolcengineAppID)
	hdr.Set("X-Api-Access-Key", cfg.Config.VolcengineAccessKey)
	hdr.Set("X-Api-Resource-Id", cfg.Config.VolcengineResourceID)
	hdr.Set("X-Api-App-Key", cfg.Config.VolcengineAppKey)
	hdr.Set("X-Api-Connect-Id", connectID)

	ws, _, err := websocket.Dial(ctx, defaultURL, &websocket.DialOptions{HTTPHeader: hdr})
	if err != nil {
		return nil, fmt.Errorf("volcengine dial: %w", err)
	}

	c := &conn{
		ws:        ws,
		sessionID: sessionID,
		instructions: "You are a helpful voice assistant.",
		voice:     "default",
	}

	if err := c.handshake(ctx); err != nil {
		_ = ws.Close(websocket.StatusInternalError, "handshake failed")
		return nil, err
	}
	return c, nil
}

type conn struct {
	ws           *websocket.Conn
	sessionID    string
	instructions string
	voice        string

	mu          sync.Mutex
	responseID  string
	audioPending bool
	readBuf     [][]byte
}

func (c *conn) handshake(ctx context.Context) error {
	startConn, err := buildFullRequest(eventStartConnection, "", map[string]any{})
	if err != nil {
		return err
	}
	if err := c.ws.Write(ctx, websocket.MessageBinary, startConn); err != nil {
		return err
	}
	if _, err := c.readUpstream(ctx); err != nil {
		return err
	}

	sessionPayload := map[string]any{
		"tts": map[string]any{
			"audio_config": map[string]any{
				"channel":     1,
				"format":      "pcm",
				"sample_rate": realtime.OutputSampleRate,
			},
		},
		"dialog": map[string]any{
			"bot_name":       "assistant",
			"system_role":    c.instructions,
			"speaking_style": "natural and concise",
		},
	}
	startSess, err := buildFullRequest(eventStartSession, c.sessionID, sessionPayload)
	if err != nil {
		return err
	}
	if err := c.ws.Write(ctx, websocket.MessageBinary, startSess); err != nil {
		return err
	}
	_, err = c.readUpstream(ctx)
	return err
}

func (c *conn) WriteClientEvent(ctx context.Context, raw []byte) error {
	evType := realtime.EventType(raw)
	switch evType {
	case realtime.EventSessionUpdate:
		var body struct {
			Session struct {
				Instructions string `json:"instructions"`
				Voice        string `json:"voice"`
			} `json:"session"`
		}
		_ = json.Unmarshal(raw, &body)
		if body.Session.Instructions != "" {
			c.instructions = body.Session.Instructions
		}
		if body.Session.Voice != "" {
			c.voice = body.Session.Voice
		}
		return nil
	case realtime.EventInputAudioBufferAppend:
		var body struct {
			Audio string `json:"audio"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}
		pcm, err := realtime.DecodePCM16Base64(body.Audio)
		if err != nil {
			return err
		}
		pcm = realtime.DownsampleTo16k(pcm)
		frame, err := buildAudioRequest(c.sessionID, pcm)
		if err != nil {
			return err
		}
		return c.ws.Write(ctx, websocket.MessageBinary, frame)
	case realtime.EventResponseCancel:
		// Best-effort interrupt: finish current session chunk handling client-side.
		c.mu.Lock()
		c.audioPending = false
		c.mu.Unlock()
		return nil
	case realtime.EventInputAudioBufferClear:
		return nil
	default:
		return nil
	}
}

func (c *conn) ReadServerEvent(ctx context.Context) ([]byte, error) {
	c.mu.Lock()
	if len(c.readBuf) > 0 {
		out := c.readBuf[0]
		c.readBuf = c.readBuf[1:]
		c.mu.Unlock()
		return out, nil
	}
	c.mu.Unlock()

	for {
		parsed, err := c.readUpstream(ctx)
		if err != nil {
			return nil, err
		}
		events := c.translate(parsed)
		if len(events) == 0 {
			continue
		}
		c.mu.Lock()
		if len(events) > 1 {
			c.readBuf = append(c.readBuf, events[1:]...)
		}
		out := events[0]
		c.mu.Unlock()
		return out, nil
	}
}

func (c *conn) readUpstream(ctx context.Context) (*parsedResponse, error) {
	typ, data, err := c.ws.Read(ctx)
	if err != nil {
		return nil, err
	}
	if typ != websocket.MessageBinary {
		return nil, fmt.Errorf("unexpected volcengine message type %d", typ)
	}
	return parseResponse(data)
}

func (c *conn) translate(p *parsedResponse) [][]byte {
	var out [][]byte

	if p.MessageType == "SERVER_ERROR" {
		return [][]byte{realtime.NewErrorEvent("upstream_error", fmt.Sprintf("volcengine error code %d", p.Code))}
	}

	// Audio payload from TTS stream
	if raw, ok := p.PayloadMsg.([]byte); ok && len(raw) > 0 {
		c.mu.Lock()
		if !c.audioPending {
			c.responseID = newID()
			c.audioPending = true
			created, _ := json.Marshal(map[string]any{
				"type":        realtime.EventResponseCreated,
				"response":    map[string]string{"id": c.responseID},
			})
			out = append(out, created)
		}
		c.mu.Unlock()

		pcm24 := realtime.UpsampleTo24k(raw)
		delta, _ := json.Marshal(map[string]any{
			"type":  realtime.EventResponseAudioDelta,
			"delta": realtime.EncodePCM16Base64(pcm24),
		})
		out = append(out, delta)
		return out
	}

	if m, ok := p.PayloadMsg.(map[string]any); ok {
		if text, ok := m["text"].(string); ok && text != "" {
			delta, _ := json.Marshal(map[string]any{
				"type":  realtime.EventResponseAudioTranscriptDelta,
				"delta": text,
			})
			out = append(out, delta)
		}
		if asr, ok := m["asr_text"].(string); ok && asr != "" {
			delta, _ := json.Marshal(map[string]any{
				"type":  realtime.EventConversationItemInputAudioTranscriptionDelta,
				"delta": asr,
			})
			out = append(out, delta)
		}
	}

	// End of utterance heuristics on ACK events
	if p.MessageType == "SERVER_ACK" && p.Event == 351 {
		out = append(out,
			mustJSON(map[string]any{"type": realtime.EventInputAudioBufferSpeechStarted}),
			mustJSON(map[string]any{"type": realtime.EventInputAudioBufferSpeechStopped}),
			mustJSON(map[string]any{"type": realtime.EventInputAudioBufferCommitted}),
		)
	}
	if p.MessageType == "SERVER_ACK" && p.Event == 359 {
		c.mu.Lock()
		if c.audioPending {
			c.audioPending = false
			c.mu.Unlock()
			out = append(out,
				mustJSON(map[string]any{"type": realtime.EventResponseAudioDone}),
				mustJSON(map[string]any{"type": realtime.EventResponseAudioTranscriptDone}),
				mustJSON(map[string]any{"type": realtime.EventResponseDone, "response": map[string]string{"id": c.responseID}}),
			)
		} else {
			c.mu.Unlock()
		}
	}
	return out
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (c *conn) Close() error {
	ctx := context.Background()
	finish, _ := buildFullRequest(eventFinishSession, c.sessionID, map[string]any{})
	_ = c.ws.Write(ctx, websocket.MessageBinary, finish)
	finishConn, _ := buildFullRequest(eventFinishConnection, "", map[string]any{})
	_ = c.ws.Write(ctx, websocket.MessageBinary, finishConn)
	return c.ws.Close(websocket.StatusNormalClosure, "closed")
}

var _ providers.UpstreamConn = (*conn)(nil)

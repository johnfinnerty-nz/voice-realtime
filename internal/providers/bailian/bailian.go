package bailian

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/lixuanqun/voice-realtime/internal/providers"
	"github.com/lixuanqun/voice-realtime/internal/realtime"
)

const defaultURL = "wss://dashscope.aliyuncs.com/api-ws/v1/inference"

// Provider adapts Alibaba Bailian multimodal dialog API to OpenAI Realtime.
type Provider struct{}

func New() *Provider { return &Provider{} }

func init() {
	providers.Register(New())
}

func (p *Provider) Name() string { return "bailian" }

func (p *Provider) Dial(ctx context.Context, cfg providers.ConnectConfig) (providers.UpstreamConn, error) {
	hdr := http.Header{}
	hdr.Set("Authorization", "Bearer "+cfg.Config.DashScopeAPIKey)

	ws, _, err := websocket.Dial(ctx, defaultURL, &websocket.DialOptions{HTTPHeader: hdr})
	if err != nil {
		return nil, fmt.Errorf("bailian dial: %w", err)
	}

	c := &conn{
		ws:          ws,
		taskID:      newID(),
		workspaceID: cfg.Config.BailianWorkspaceID,
		appID:       cfg.Config.BailianAppID,
		userID:      cfg.Config.BailianUserID,
		voice:       "longxiaochun_v2",
		instructions: "",
		state:       stateInit,
	}
	if err := c.startSession(ctx); err != nil {
		_ = ws.Close(websocket.StatusInternalError, "start failed")
		return nil, err
	}
	return c, nil
}

type dialogState int

const (
	stateInit dialogState = iota
	stateStarted
	stateListening
	stateResponding
)

type conn struct {
	ws          *websocket.Conn
	taskID      string
	dialogID    string
	workspaceID string
	appID       string
	userID      string
	voice       string
	instructions string

	mu          sync.Mutex
	state       dialogState
	canSendAudio bool
	responseID  string
	readBuf     [][]byte
}

func (c *conn) startSession(ctx context.Context) error {
	msg := map[string]any{
		"header": map[string]any{
			"action":    "run-task",
			"task_id":   c.taskID,
			"streaming": "duplex",
		},
		"payload": map[string]any{
			"task_group": "aigc",
			"task":       "multimodal-generation",
			"function":   "generation",
			"model":      "multimodal-dialog",
			"input": map[string]any{
				"directive":     "Start",
				"workspace_id": c.workspaceID,
				"app_id":        c.appID,
			},
			"parameters": map[string]any{
				"upstream": map[string]any{
					"type":        "AudioOnly",
					"mode":        "duplex",
					"sample_rate": realtime.VolcSampleRate,
				},
				"downstream": map[string]any{
					"voice":       c.voice,
					"sample_rate": realtime.OutputSampleRate,
				},
				"client_info": map[string]any{
					"user_id": c.userID,
				},
			},
		},
	}
	raw, _ := json.Marshal(msg)
	if err := c.ws.Write(ctx, websocket.MessageText, raw); err != nil {
		return err
	}

	// Wait until Listening before accepting audio.
	for c.state != stateListening {
		typ, data, err := c.ws.Read(ctx)
		if err != nil {
			return err
		}
		var events [][]byte
		if typ == websocket.MessageBinary {
			events = c.translateBinary(data)
		} else {
			events = c.translateText(data)
		}
		_ = events
	}
	return nil
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
		c.mu.Lock()
		if body.Session.Voice != "" {
			c.voice = body.Session.Voice
		}
		if body.Session.Instructions != "" {
			c.instructions = body.Session.Instructions
		}
		c.mu.Unlock()
		return nil
	case realtime.EventInputAudioBufferAppend:
		c.mu.Lock()
		ok := c.canSendAudio
		c.mu.Unlock()
		if !ok {
			return nil
		}
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
		return c.ws.Write(ctx, websocket.MessageBinary, pcm)
	case realtime.EventResponseCancel:
		return c.sendDirective(ctx, "RequestToSpeak")
	case realtime.EventInputAudioBufferClear:
		return c.sendDirective(ctx, "CancelSpeech")
	default:
		return nil
	}
}

func (c *conn) sendDirective(ctx context.Context, directive string) error {
	c.mu.Lock()
	dialogID := c.dialogID
	c.mu.Unlock()
	msg := map[string]any{
		"header": map[string]any{
			"action":    "continue-task",
			"task_id":   c.taskID,
			"streaming": "duplex",
		},
		"payload": map[string]any{
			"input": map[string]any{
				"directive": directive,
				"dialog_id": dialogID,
			},
		},
	}
	raw, _ := json.Marshal(msg)
	return c.ws.Write(ctx, websocket.MessageText, raw)
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
		typ, data, err := c.ws.Read(ctx)
		if err != nil {
			return nil, err
		}
		var events [][]byte
		if typ == websocket.MessageBinary {
			events = c.translateBinary(data)
		} else {
			events = c.translateText(data)
		}
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

func (c *conn) translateText(raw []byte) [][]byte {
	var msg struct {
		Header struct {
			Event string `json:"event"`
		} `json:"header"`
		Payload struct {
			Output struct {
				Event    string `json:"event"`
				State    string `json:"state"`
				DialogID string `json:"dialog_id"`
				Text     string `json:"text"`
				Type     string `json:"type"`
			} `json:"output"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil
	}

	ev := msg.Payload.Output.Event
	if msg.Payload.Output.DialogID != "" {
		c.mu.Lock()
		c.dialogID = msg.Payload.Output.DialogID
		c.mu.Unlock()
	}

	var out [][]byte

	switch ev {
	case "Started":
		c.mu.Lock()
		c.state = stateStarted
		c.mu.Unlock()
		out = append(out, realtime.NewSessionCreated(c.taskID))
	case "DialogStateChanged":
		c.mu.Lock()
		switch msg.Payload.Output.State {
		case "Listening":
			c.state = stateListening
			c.canSendAudio = true
		case "Responding":
			c.state = stateResponding
			c.canSendAudio = false
		}
		c.mu.Unlock()
		updated, _ := json.Marshal(map[string]any{
			"type": realtime.EventSessionUpdated,
			"session": map[string]any{
				"dialog_state": msg.Payload.Output.State,
			},
		})
		out = append(out, updated)
	case "SpeechStarted":
		out = append(out, mustJSON(map[string]any{"type": realtime.EventInputAudioBufferSpeechStarted}))
	case "SpeechEnded":
		out = append(out,
			mustJSON(map[string]any{"type": realtime.EventInputAudioBufferSpeechStopped}),
			mustJSON(map[string]any{"type": realtime.EventInputAudioBufferCommitted}),
		)
	case "RespondingStarted":
		c.mu.Lock()
		c.responseID = newID()
		id := c.responseID
		c.mu.Unlock()
		out = append(out, mustJSON(map[string]any{
			"type":     realtime.EventResponseCreated,
			"response": map[string]string{"id": id},
		}))
	case "RespondingEnded":
		c.mu.Lock()
		id := c.responseID
		c.mu.Unlock()
		out = append(out,
			mustJSON(map[string]any{"type": realtime.EventResponseAudioDone}),
			mustJSON(map[string]any{"type": realtime.EventResponseDone, "response": map[string]string{"id": id}}),
		)
	case "RespondingContent":
		if msg.Payload.Output.Text != "" {
			if msg.Payload.Output.Type == "transcript" {
				out = append(out, mustJSON(map[string]any{
					"type":  realtime.EventConversationItemInputAudioTranscriptionDelta,
					"delta": msg.Payload.Output.Text,
				}))
			} else {
				out = append(out, mustJSON(map[string]any{
					"type":  realtime.EventResponseAudioTranscriptDelta,
					"delta": msg.Payload.Output.Text,
				}))
			}
		}
	}
	return out
}

func (c *conn) translateBinary(pcm []byte) [][]byte {
	if len(pcm) == 0 {
		return nil
	}
	pcm24 := realtime.UpsampleTo24k(pcm)
	delta, _ := json.Marshal(map[string]any{
		"type":  realtime.EventResponseAudioDelta,
		"delta": realtime.EncodePCM16Base64(pcm24),
	})
	return [][]byte{delta}
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func (c *conn) Close() error {
	ctx := context.Background()
	_ = c.sendDirective(ctx, "Stop")
	return c.ws.Close(websocket.StatusNormalClosure, "closed")
}

var _ providers.UpstreamConn = (*conn)(nil)

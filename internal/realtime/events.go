package realtime

import "encoding/json"

// Common OpenAI Realtime API event type constants.
const (
	EventSessionUpdate              = "session.update"
	EventSessionCreated             = "session.created"
	EventSessionUpdated             = "session.updated"
	EventInputAudioBufferAppend     = "input_audio_buffer.append"
	EventInputAudioBufferClear      = "input_audio_buffer.clear"
	EventInputAudioBufferCommit     = "input_audio_buffer.commit"
	EventInputAudioBufferSpeechStarted = "input_audio_buffer.speech_started"
	EventInputAudioBufferSpeechStopped = "input_audio_buffer.speech_stopped"
	EventInputAudioBufferCommitted  = "input_audio_buffer.committed"
	EventResponseCreate             = "response.create"
	EventResponseCancel             = "response.cancel"
	EventResponseCreated            = "response.created"
	EventResponseAudioDelta         = "response.audio.delta"
	EventResponseAudioDone          = "response.audio.done"
	EventResponseAudioTranscriptDelta = "response.audio_transcript.delta"
	EventResponseAudioTranscriptDone  = "response.audio_transcript.done"
	EventResponseDone               = "response.done"
	EventConversationItemInputAudioTranscriptionDelta = "conversation.item.input_audio_transcription.delta"
	EventConversationItemInputAudioTranscriptionCompleted = "conversation.item.input_audio_transcription.completed"
	EventError                      = "error"
)

// ClientEvent is a minimal parsed client event.
type ClientEvent struct {
	Type  string          `json:"type"`
	Audio string          `json:"audio,omitempty"`
	Session json.RawMessage `json:"session,omitempty"`
}

// ParseClientEvent unmarshals a client JSON event.
func ParseClientEvent(raw []byte) (*ClientEvent, error) {
	var ev ClientEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}

// EventType extracts the type field without full unmarshaling.
func EventType(raw []byte) string {
	var partial struct {
		Type string `json:"type"`
	}
	_ = json.Unmarshal(raw, &partial)
	return partial.Type
}

// NewErrorEvent builds an OpenAI-style error event.
func NewErrorEvent(code, message string) []byte {
	b, _ := json.Marshal(map[string]any{
		"type": EventError,
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
	return b
}

// NewSessionCreated returns a session.created event.
func NewSessionCreated(sessionID string) []byte {
	b, _ := json.Marshal(map[string]any{
		"type": EventSessionCreated,
		"session": map[string]any{
			"id": sessionID,
		},
	})
	return b
}

package volcengine

import (
	"testing"
)

func TestParseResponseEmpty(t *testing.T) {
	_, err := parseResponse([]byte{0x10, 0x10, 0x11, 0x00})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildAudioRequest(t *testing.T) {
	frame, err := buildAudioRequest("sess-1", []byte{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(frame) < 8 {
		t.Fatalf("frame too short: %d", len(frame))
	}
}

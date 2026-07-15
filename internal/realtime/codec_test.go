package realtime

import (
	"encoding/binary"
	"testing"
)

func TestResamplePCM16RoundTripLength(t *testing.T) {
	in := make([]byte, 24000*2) // 1 second at 24kHz
	for i := 0; i < 24000; i++ {
		binary.LittleEndian.PutUint16(in[i*2:], uint16(i%1000))
	}
	down := DownsampleTo16k(in)
	if len(down) != 16000*2 {
		t.Fatalf("expected 16000 samples, got %d bytes", len(down))
	}
	up := UpsampleTo24k(down)
	if len(up) != 24000*2 {
		t.Fatalf("expected 24000 samples, got %d bytes", len(up))
	}
}

func TestEncodeDecodeBase64(t *testing.T) {
	raw := []byte{0x01, 0x02, 0x03, 0x04}
	b64 := EncodePCM16Base64(raw)
	out, err := DecodePCM16Base64(b64)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(raw) {
		t.Fatalf("mismatch: %v vs %v", out, raw)
	}
}

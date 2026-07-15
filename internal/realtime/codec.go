package realtime

import (
	"encoding/base64"
	"encoding/binary"
)

const (
	InputSampleRate  = 24000 // OpenAI Realtime client default
	VolcSampleRate   = 16000 // Volcengine / Bailian upstream default
	OutputSampleRate = 24000
)

// DecodePCM16Base64 decodes base64 PCM16 little-endian audio.
func DecodePCM16Base64(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)
}

// EncodePCM16Base64 encodes PCM16 bytes to base64.
func EncodePCM16Base64(pcm []byte) string {
	return base64.StdEncoding.EncodeToString(pcm)
}

// ResamplePCM16 resamples mono PCM16 between sample rates using linear interpolation.
func ResamplePCM16(pcm []byte, fromRate, toRate int) []byte {
	if fromRate == toRate || len(pcm) < 2 {
		return pcm
	}
	samples := len(pcm) / 2
	outSamples := samples * toRate / fromRate
	if outSamples == 0 {
		return nil
	}
	out := make([]byte, outSamples*2)
	for i := 0; i < outSamples; i++ {
		srcPos := float64(i) * float64(fromRate) / float64(toRate)
		idx := int(srcPos)
		if idx >= samples-1 {
			idx = samples - 2
		}
		frac := srcPos - float64(idx)
		s0 := int16(binary.LittleEndian.Uint16(pcm[idx*2:]))
		s1 := int16(binary.LittleEndian.Uint16(pcm[(idx+1)*2:]))
		v := int16(float64(s0)*(1-frac) + float64(s1)*frac)
		binary.LittleEndian.PutUint16(out[i*2:], uint16(v))
	}
	return out
}

// DownsampleTo16k converts 24kHz PCM16 to 16kHz for upstream providers.
func DownsampleTo16k(pcm24 []byte) []byte {
	return ResamplePCM16(pcm24, InputSampleRate, VolcSampleRate)
}

// UpsampleTo24k converts 16kHz PCM16 to 24kHz for OpenAI clients.
func UpsampleTo24k(pcm16 []byte) []byte {
	return ResamplePCM16(pcm16, VolcSampleRate, OutputSampleRate)
}

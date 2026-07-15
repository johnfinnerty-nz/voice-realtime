package volcengine

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func generateHeader(messageType, serialMethod, compression int) []byte {
	header := make([]byte, 4)
	header[0] = (protocolVersion << 4) | 0x01
	header[1] = (byte(messageType) << 4) | byte(msgWithEvent)
	header[2] = (byte(serialMethod) << 4) | byte(compression)
	header[3] = 0x00
	return header
}

func buildFullRequest(event int, sessionID string, payload any) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	compressed := gzipCompress(payloadBytes)

	buf := bytes.NewBuffer(nil)
	buf.Write(generateHeader(clientFullRequest, jsonSerialization, gzipCompression))
	_ = binary.Write(buf, binary.BigEndian, int32(event))
	if sessionID != "" {
		sid := []byte(sessionID)
		_ = binary.Write(buf, binary.BigEndian, int32(len(sid)))
		buf.Write(sid)
	}
	_ = binary.Write(buf, binary.BigEndian, int32(len(compressed)))
	buf.Write(compressed)
	return buf.Bytes(), nil
}

func buildAudioRequest(sessionID string, audio []byte) ([]byte, error) {
	compressed := gzipCompress(audio)
	buf := bytes.NewBuffer(nil)
	buf.Write(generateHeader(clientAudioOnlyRequest, noSerialization, gzipCompression))
	_ = binary.Write(buf, binary.BigEndian, int32(eventTaskRequest))
	sid := []byte(sessionID)
	_ = binary.Write(buf, binary.BigEndian, int32(len(sid)))
	buf.Write(sid)
	_ = binary.Write(buf, binary.BigEndian, int32(len(compressed)))
	buf.Write(compressed)
	return buf.Bytes(), nil
}

func gzipCompress(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, _ = w.Write(data)
	_ = w.Close()
	return b.Bytes()
}

func gzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	_, err = out.ReadFrom(r)
	return out.Bytes(), err
}

type parsedResponse struct {
	MessageType string
	Event       int
	SessionID   string
	PayloadMsg  any
	PayloadRaw  []byte
	Code        int
}

func parseResponse(res []byte) (*parsedResponse, error) {
	if len(res) < 4 {
		return nil, fmt.Errorf("response too short")
	}
	headerSize := int(res[0] & 0x0f)
	messageType := res[1] >> 4
	flags := res[1] & 0x0f
	serialMethod := res[2] >> 4
	compression := res[2] & 0x0f

	payload := res[headerSize*4:]
	result := &parsedResponse{}

	switch messageType {
	case serverFullResponse:
		result.MessageType = "SERVER_FULL_RESPONSE"
	case serverACK:
		result.MessageType = "SERVER_ACK"
	case serverErrorResponse:
		result.MessageType = "SERVER_ERROR"
	default:
		result.MessageType = "UNKNOWN"
	}

	start := 0
	if flags&msgWithEvent > 0 && len(payload) >= 4 {
		result.Event = int(binary.BigEndian.Uint32(payload[:4]))
		start = 4
	}
	payload = payload[start:]

	if messageType == serverErrorResponse {
		if len(payload) >= 8 {
			result.Code = int(binary.BigEndian.Uint32(payload[:4]))
			size := int(binary.BigEndian.Uint32(payload[4:8]))
			raw := payload[8:]
			if len(raw) > size {
				raw = raw[:size]
			}
			result.PayloadRaw = raw
		}
		return result, nil
	}

	if len(payload) < 4 {
		return result, nil
	}
	sidLen := int(binary.BigEndian.Uint32(payload[:4]))
	if sidLen < 0 || len(payload) < 4+sidLen+4 {
		return result, nil
	}
	result.SessionID = string(payload[4 : 4+sidLen])
	rest := payload[4+sidLen:]
	if len(rest) < 4 {
		return result, nil
	}
	dataLen := int(binary.BigEndian.Uint32(rest[:4]))
	raw := rest[4:]
	if len(raw) > dataLen {
		raw = raw[:dataLen]
	}

	if compression == gzipCompression {
		decompressed, err := gzipDecompress(raw)
		if err != nil {
			return nil, err
		}
		raw = decompressed
	}

	result.PayloadRaw = raw
	if serialMethod == jsonSerialization && len(raw) > 0 {
		var v any
		if err := json.Unmarshal(raw, &v); err == nil {
			result.PayloadMsg = v
		}
	} else if serialMethod == noSerialization {
		result.PayloadMsg = raw
	}
	return result, nil
}

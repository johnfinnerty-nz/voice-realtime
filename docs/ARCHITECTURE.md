# Architecture

## Overview

`voice-realtime` is an OpenAI Realtime-compatible WebSocket gateway that routes client sessions to cloud end-to-end voice LLM backends.

```
Client (OpenAI Realtime protocol)
        │
        ▼
┌───────────────────┐
│  Gateway          │  ws://host/v1/realtime?provider=&model=
│  /v1/realtime     │
└─────────┬─────────┘
          │
    ┌─────┴─────┬──────────┬──────────┐
    ▼           ▼          ▼          ▼
 zhipu      stepfun   volcengine   bailian
 (proxy)    (proxy)   (translate)  (translate)
```

## Northbound API

- **Endpoint**: `GET /v1/realtime` (WebSocket upgrade)
- **Query**: `provider` (required for non-default), `model` (optional)
- **Audio**: PCM16 little-endian, 24 kHz mono, base64 in JSON events
- **Events**: subset of OpenAI Realtime — see [Vui realtime-api](https://github.com/lixuanqun/vui/blob/main/docs/realtime-api.md)

### Session lifecycle

1. Client connects with `provider` and `model` query params.
2. Gateway validates credentials for the provider, dials upstream.
3. Two goroutines relay:
   - client → `UpstreamConn.WriteClientEvent`
   - `UpstreamConn.ReadServerEvent` → client
4. `response.cancel` sets session cancel flag; subsequent `response.audio.delta` from upstream is dropped until `response.done`.

## Provider plugin model

```go
type Provider interface {
    Name() string
    Dial(ctx context.Context, cfg ConnectConfig) (UpstreamConn, error)
}

type UpstreamConn interface {
    WriteClientEvent(ctx context.Context, raw []byte) error
    ReadServerEvent(ctx context.Context) ([]byte, error)
    Close() error
}
```

Providers register via `init()` + blank import in `cmd/voice-realtime/main.go`.

## Adapter strategies

| Provider | Strategy | Notes |
|----------|----------|-------|
| zhipu | Transparent WebSocket proxy | Injects `Authorization: Bearer` |
| stepfun | Transparent WebSocket proxy | Adds `?model=` query param |
| volcengine | Binary protocol translation | Gzip-framed events; 16 kHz upstream audio |
| bailian | JSON + binary state machine | Waits for `DialogStateChanged: Listening` before audio |

## Audio pipeline

- **Client → gateway**: 24 kHz PCM16 base64 (OpenAI default)
- **Gateway → volcengine/bailian**: resample to 16 kHz via linear interpolation (`internal/realtime/codec.go`)
- **Upstream → client**: upsample 16 kHz → 24 kHz for TTS chunks

## Security

- Cloud API keys live only on the server (environment variables).
- Do not expose the gateway to the public internet without an auth layer (TODO).

## Observability

- Structured JSON logs via `log/slog`
- `GET /healthz` for liveness

## Adding a provider

1. Create `internal/providers/<name>/`
2. Implement `Provider` and `UpstreamConn`
3. Call `providers.Register` in `init()`
4. Add blank import in `cmd/voice-realtime/main.go`
5. Add env validation in `internal/config/config.go`
6. Document in `docs/providers/<name>.md`

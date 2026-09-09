# Python realtime client

This minimal client connects to a local `voice-realtime` gateway using the OpenAI Realtime WebSocket event format. It sends `session.update`, optionally appends a raw PCM16 audio file with `input_audio_buffer.append`, and saves received `response.audio.delta` audio to a PCM16 file.

## Prerequisites

- Python 3.10 or later
- A running local gateway. From the repository root, configure a provider key in `.env` and run `go run ./cmd/voice-realtime`.

## Run

```bash
cd examples/python_client
python -m venv .venv
.venv\Scripts\activate
pip install -r requirements.txt
python realtime_client.py
```

The default URL uses the Zhipu provider. Change it with `--url` to select another supported provider:

```bash
python realtime_client.py --url "ws://localhost:8080/v1/realtime?provider=stepfun&model=stepaudio-2.5-realtime"
```

To send audio, provide a raw 24 kHz mono PCM16 file. The returned PCM16 data is saved to `response.pcm` by default:

```bash
python realtime_client.py --audio input.pcm --output response.pcm
```

The example deliberately keeps playback out of scope. Use any PCM16-capable player or convert the output to WAV for local playback.

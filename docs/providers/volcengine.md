# 火山引擎豆包实时对话

## 上游

- URL: `wss://openspeech.bytedance.com/api/v3/realtime/dialogue`
- 鉴权 Header:
  - `X-Api-App-ID`
  - `X-Api-Access-Key`
  - `X-Api-Resource-Id`（默认 `volc.speech.dialog`）
  - `X-Api-App-Key`
  - `X-Api-Connect-Id`（每次连接 UUID）

## 网关连接

```
ws://localhost:8080/v1/realtime?provider=volcengine
```

## 环境变量

```bash
export VOLCENGINE_APP_ID=...
export VOLCENGINE_ACCESS_KEY=...
```

## 适配方式

火山使用 **gzip 压缩的二进制帧协议**（非 JSON WebSocket）。网关负责：

| OpenAI Realtime | 火山 |
|-----------------|------|
| `session.update` → `instructions` | 写入 `dialog.system_role`（建连时） |
| `input_audio_buffer.append` | event 200 音频帧（16 kHz PCM） |
| `response.audio.delta` | TTS 二进制流 → base64 24 kHz |
| `response.cancel` | 客户端侧丢弃进行中的音频 |

## 参考

- [RealtimeDialog-doubao 示例](https://github.com/SUAT-AIRI/RealtimeDialog-doubao)
- [豆包实时语音技术解析](https://www.volcengine.com/docs/6360/1330208)

# 阿里云百炼多模态交互

## 上游

- URL: `wss://dashscope.aliyuncs.com/api-ws/v1/inference`
- 鉴权: `Authorization: Bearer <DASHSCOPE_API_KEY>`

## 网关连接

```
ws://localhost:8080/v1/realtime?provider=bailian
```

## 环境变量

```bash
export DASHSCOPE_API_KEY=...
export BAILIAN_WORKSPACE_ID=llm-xxxxxxxx
export BAILIAN_APP_ID=xxxxxxxx
```

在百炼控制台「多模态交互开发套件」创建应用后获取 `workspace_id` 与 `app_id`。

## 状态机

百炼协议与 OpenAI Realtime 差异较大，网关维护以下映射：

```
Start → Started → DialogStateChanged(Listening) → [可发音频]
```

| 百炼事件 | OpenAI Realtime |
|----------|-----------------|
| `SpeechStarted` | `input_audio_buffer.speech_started` |
| `SpeechEnded` | `input_audio_buffer.speech_stopped` + `committed` |
| `RespondingStarted` | `response.created` |
| 二进制音频 | `response.audio.delta` |
| `RespondingEnded` | `response.audio.done` + `response.done` |
| `RespondingContent` (dialog) | `response.audio_transcript.delta` |

| OpenAI Realtime | 百炼 |
|-----------------|------|
| `input_audio_buffer.append` | 二进制 PCM 16 kHz |
| `response.cancel` | `RequestToSpeak` |
| `input_audio_buffer.clear` | `CancelSpeech` |

## 模式

默认使用 `duplex` 双工模式（服务端 VAD），对应 OpenAI `turn_detection: server_vad`。

## 参考

- [实时多模态交互 API](https://help.aliyun.com/zh/model-studio/multimodal-interaction-protocol)

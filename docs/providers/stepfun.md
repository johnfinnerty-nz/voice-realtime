# 阶跃星辰 StepFun Realtime

## 上游

- URL: `wss://api.stepfun.com/v1/realtime?model=<model>`
- 鉴权: `Authorization: Bearer <STEPFUN_API_KEY>`

## 网关连接

```
ws://localhost:8080/v1/realtime?provider=stepfun&model=stepaudio-2.5-realtime
```

## 模型

| model | 说明 |
|-------|------|
| `stepaudio-2.5-realtime` | 端到端实时语音（默认） |
| `step-audio-2` | 上一代，支持音色复刻 |
| `step-audio-2-mini` | 轻量版 |

## 适配方式

透明代理，协议与 OpenAI Realtime 一致。

## 注意

- 会话最长约 30 分钟
- 模型首次以音频响应后，`voice` 不可再修改

## 参考

- [双向实时语音 API](https://platform.stepfun.com/docs/zh/api-reference/realtime/chat)
- [实时对话开发指南](https://platform.stepfun.com/docs/zh/guides/developer/realtime)

# 智谱 GLM-Realtime

## 上游

- URL: `wss://open.bigmodel.cn/api/paas/v4/realtime`
- 鉴权: `Authorization: Bearer <ZHIPU_API_KEY>`

## 网关连接

```
ws://localhost:8080/v1/realtime?provider=zhipu&model=glm-realtime-flash
```

## 模型

| model | 说明 |
|-------|------|
| `glm-realtime-flash` | 9B 轻量（默认） |
| `glm-realtime-air` | 32B |

## 适配方式

透明代理：北向 OpenAI Realtime 事件与智谱上游 1:1 透传。

## 参考

- [智谱 GLM-Realtime 文档](https://docs.bigmodel.cn/cn/guide/models/sound-and-video/glm-realtime)
- [glm-realtime-sdk](https://github.com/MetaGLM/glm-realtime-sdk)

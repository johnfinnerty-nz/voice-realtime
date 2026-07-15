# voice-realtime

[![CI](https://github.com/lixuanqun/voice-realtime/actions/workflows/ci.yml/badge.svg)](https://github.com/lixuanqun/voice-realtime/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev/dl/)

**多云端到端语音大模型网关** — 对外暴露 [OpenAI Realtime API](https://platform.openai.com/docs/guides/realtime) 兼容 WebSocket，统一接入火山豆包、阿里云百炼、智谱 GLM-Realtime、阶跃星辰 StepFun Realtime。

```bash
export ZHIPU_API_KEY=your_key
go run ./cmd/voice-realtime
# ws://localhost:8080/v1/realtime?provider=zhipu&model=glm-realtime-flash
```

## 支持的云厂商

| Provider | 查询参数 `provider` | 默认模型 | 适配方式 |
|----------|---------------------|----------|----------|
| 智谱 | `zhipu` | `glm-realtime-flash` | OpenAI Realtime 透明代理 |
| 阶跃星辰 | `stepfun` | `stepaudio-2.5-realtime` | OpenAI Realtime 透明代理 |
| 火山豆包 | `volcengine` | — | 二进制协议双向转换 |
| 阿里云百炼 | `bailian` | `multimodal-dialog` | 多模态交互 API 状态机 |

## 快速开始

### 1. 配置环境变量

复制 [`.env.example`](.env.example) 并填入对应厂商密钥（只需配置你要用的厂商）：

```bash
cp .env.example .env
```

### 2. 启动网关

```bash
go run ./cmd/voice-realtime
```

### 3. 连接

```
ws://localhost:8080/v1/realtime?provider=stepfun&model=stepaudio-2.5-realtime
```

客户端使用 OpenAI Realtime 协议：`session.update`、`input_audio_buffer.append`（PCM16 24kHz base64）、`response.audio.delta` 等。详见 [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)。

## 环境变量

| 变量 | 用途 |
|------|------|
| `VOICE_REALTIME_ADDR` | 监听地址，默认 `:8080` |
| `ZHIPU_API_KEY` | 智谱 API Key |
| `STEPFUN_API_KEY` | 阶跃星辰 API Key |
| `VOLCENGINE_APP_ID` / `VOLCENGINE_ACCESS_KEY` | 火山豆包实时对话 |
| `DASHSCOPE_API_KEY` | 阿里云百炼 |
| `BAILIAN_WORKSPACE_ID` / `BAILIAN_APP_ID` | 百炼多模态应用 |

## 项目结构

```
cmd/voice-realtime/     # 入口
internal/gateway/       # WebSocket 网关
internal/realtime/      # OpenAI Realtime 事件与音频编解码
internal/providers/     # 云厂商插件（zhipu / stepfun / volcengine / bailian）
docs/                   # 架构与厂商接入文档
```

## 开发

```bash
make test    # 运行测试
make build   # 编译二进制
make run     # 启动服务
```

## 相关项目

本仓库位于 [voice_repo](https://github.com/lixuanqun) monorepo 的 `realtime/` 目录，与 [Vui](https://github.com/lixuanqun) Realtime API 规范对齐，可与 OpenClaw、自研客户端对接。

## License

MIT

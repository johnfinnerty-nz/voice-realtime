# 参与贡献

感谢你对 **voice-realtime** 的关注！本项目旨在降低国内多云端到端语音大模型的接入门槛，非常欢迎各类贡献。

---

## 我能贡献什么？

| 类型 | 示例 | 适合人群 |
|------|------|----------|
| **代码** | 新 Provider、协议映射、测试 | Go 开发者 |
| **联调** | 用真实云账号验证对话流程 | 有火山/百炼/智谱/阶跃账号的同学 |
| **文档** | 协议映射表、教程、英文翻译 | 技术写作者 |
| **示例** | Python/JS/Web 最小客户端 | 前端 / 全栈 |
| **Issue** | Bug 报告、功能建议、RFC | 所有人 |

---

## 开发环境

```bash
git clone https://github.com/lixuanqun/voice-realtime.git
cd voice-realtime
cp .env.example .env   # 填入你要测的厂商 Key

go test ./...
go run ./cmd/voice-realtime
```

要求：**Go 1.21+**

---

## 贡献流程

1. **先看路线** — [docs/ROADMAP.md](docs/ROADMAP.md) 确认任务在哪个 Phase
2. **搜 Issue** — 避免重复劳动；没有合适 Issue 就先开一个讨论
3. **Fork & Branch** — `git checkout -b feat/provider-xxx` 或 `fix/bailian-listening`
4. **开发与测试** — `go test ./...` 必须通过
5. **提交 PR** — 填写 PR 模板，关联 `Closes #123`

### Commit 规范

```
<type>(<scope>): <subject>

feat(volcengine): map TTS audio to response.audio.delta
fix(bailian): wait for Listening before sending audio
docs(readme): add architecture diagram link
test(codec): add 24k to 16k resample cases
```

`type`: `feat` | `fix` | `docs` | `test` | `refactor` | `ci`

---

## 新增云厂商 Provider

这是最有价值的贡献之一。步骤：

### 1. 调研协议

- 上游 WebSocket URL、鉴权方式
- 与 OpenAI Realtime 的差异点（事件名、音频格式、状态机）
- 在 `docs/providers/<name>.md` 写 **映射表**

### 2. 选择适配策略

| 策略 | 何时用 |
|------|--------|
| **Proxy**（`internal/providers/proxy`） | 上游已是 OpenAI Realtime 风格 |
| **Translate**（自实现 `UpstreamConn`） | 二进制协议或自定义 JSON 状态机 |

### 3. 实现代码

```
internal/providers/<name>/
  <name>.go       # Provider + UpstreamConn
  <name>_test.go  # 契约测试（mock upstream）
```

```go
func init() {
    providers.Register(New())
}
```

并在 `cmd/voice-realtime/main.go` 添加 blank import：

```go
_ "github.com/lixuanqun/voice-realtime/internal/providers/<name>"
```

### 4. 配置

- `internal/config/config.go` — 环境变量校验
- `.env.example` — 示例

### 5. PR Checklist

- [ ] `go test ./...` 通过
- [ ] `docs/providers/<name>.md` 映射表完整
- [ ] `.env.example` 已更新
- [ ] README 厂商表已更新（如需要）

---

## Issue 与 Label

提交 Issue 时请选择合适的 Label，详见 [.github/LABELS.md](.github/LABELS.md)。

| 场景 | 推荐 Label |
|------|------------|
| 想接新手任务 | `good first issue` |
| 需要联调帮助 | `help wanted` |
| 新厂商请求 | `provider/new` + `enhancement` |
| 架构讨论 | `type/rfc` |

---

## 行为准则

- 尊重他人，技术讨论对事不对人
- 不提交 API Key、`.env` 等敏感信息
- 大型改动（>500 行或跨 Phase）请先开 RFC Issue

---

## 联系

- [GitHub Issues](https://github.com/lixuanqun/voice-realtime/issues)
- [GitHub Discussions](https://github.com/lixuanqun/voice-realtime/discussions)（启用后）

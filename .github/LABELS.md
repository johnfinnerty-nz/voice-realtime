# GitHub Labels 说明

本仓库使用分层 Label，方便贡献者快速筛选 Issue 和 PR。

## 类型 Type

| Label | 颜色 | 用途 |
|-------|------|------|
| `bug` | 红 | 功能不正常 |
| `enhancement` | 蓝 | 新功能或改进 |
| `documentation` | 深蓝 | 文档相关 |
| `question` | 紫 | 需要更多信息 |
| `type/rfc` | 紫 | 架构/协议设计讨论，合并前需达成共识 |
| `duplicate` | 灰 | 重复 Issue |
| `invalid` | 黄 | 无效 Issue |
| `wontfix` | 白 | 不计划修复 |

## 领域 Area

| Label | 用途 |
|-------|------|
| `area/gateway` | WebSocket 接入、路由、会话、鉴权 |
| `area/protocol` | OpenAI Realtime 事件映射、音频编解码 |
| `area/provider` | Provider 插件框架、注册表 |
| `area/observability` | 日志、metrics、trace |
| `area/docs` | README、架构文档、教程 |
| `area/ci` | GitHub Actions、测试基础设施 |

## 云厂商 Provider

| Label | 厂商 |
|-------|------|
| `provider/zhipu` | 智谱 GLM-Realtime |
| `provider/stepfun` | 阶跃星辰 StepFun |
| `provider/volcengine` | 火山豆包 |
| `provider/bailian` | 阿里云百炼 |
| `provider/new` | 请求接入新厂商 |

## 开发阶段 Roadmap

| Label | Phase | 说明 |
|-------|-------|------|
| `roadmap/phase-0` | P0 | 脚手架、文档、CI |
| `roadmap/phase-1` | P1 | 智谱 + 阶跃透明代理 |
| `roadmap/phase-2` | P2 | 火山协议转换 |
| `roadmap/phase-3` | P3 | 百炼状态机 |
| `roadmap/phase-4` | P4 | 生产加固 |

## 优先级 Priority

| Label | 说明 |
|-------|------|
| `priority/high` | 阻塞主线或影响联调 |
| `priority/medium` | 正常排期 |
| `priority/low` | 锦上添花 |

## 贡献者友好

| Label | 说明 |
|-------|------|
| `good first issue` | 适合新手，有清晰范围 |
| `help wanted` | 需要社区帮助（尤其联调、文档） |

## 状态 Status

| Label | 说明 |
|-------|------|
| `status/blocked` | 依赖外部条件（账号、文档、上游变更） |
| `status/in-progress` | 已有人认领开发中 |

## Issue 打标建议

```text
Bug 报告:        bug + area/* + provider/*
功能请求:        enhancement + area/* + roadmap/*
新厂商:          enhancement + provider/new + type/rfc
联调求助:        help wanted + provider/*
文档:            documentation + area/docs
```

## 创建 Labels

维护者可用脚本批量创建：

```bash
bash .github/scripts/create-labels.sh
```

或手动：

```bash
gh label create "area/gateway" --color "1d76db" --description "WebSocket gateway, routing, session" -R lixuanqun/voice-realtime
```

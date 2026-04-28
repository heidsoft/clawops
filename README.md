# ClawOps - AI DevOps Agent Platform

> **你的第二个 DevOps 工程师**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/heidsoft/clawops)](https://github.com/heidsoft/clawops/stargazers)

## 🎯 是什么

ClawOps 是一个**开源的 AI 驱动的 DevOps 自动化平台**。用自然语言就能完成部署、监控、备份等运维工作。

```
用户: "帮我部署一个 nginx 到阿里云"
AI:   "我来帮你完成部署，请确认以下配置..."
       → 创建实例 → 配置安全组 → 安装 Docker → 部署应用
AI:   "✅ 部署完成！访问地址: http://101.200.XXX.XXX"
```

## ✨ 特性

- 🗣️ **自然语言部署** - 用说话的方式完成部署
- 🤖 **AI 驱动** - LLM 理解意图，自动执行
- 🔌 **Skill 可扩展** - 像搭积木一样扩展能力
- ☁️ **多云支持** - 阿里云、AWS、Docker、K8s
- 🔔 **智能告警** - 异常自动告警 + AI 分析
- 🏢 **企业级** - 多租户、权限、审计

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- (可选) Ollama 或 OpenAI API

### 安装

```bash
# 克隆仓库
git clone https://github.com/heidsoft/clawops.git
cd clawops

# 启动后端
cd backend
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml 配置数据库和云 API
go run cmd/server/main.go

# 启动前端 (新窗口)
cd frontend
npm install
npm run dev
```

### 配置

编辑 `backend/config/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your-password
  dbname: clawops

llm:
  provider: openai  # 或 ollama
  api_key: your-api-key
  base_url: https://api.openai.com/v1  # ollama 用 http://localhost:11434

cloud:
  aliyun:
    access_key: your-access-key
    secret_key: your-secret-key
    region: cn-beijing

dingtalk:
  webhook: your-dingtalk-webhook
```

### 使用

1. 打开浏览器 http://localhost:3000
2. 登录后进入 AI 对话界面
3. 输入自然语言指令：
   - "部署一个 nginx"
   - "查看当前所有部署"
   - "服务器 CPU 告警了，帮我分析一下"
   - "备份数据库"

## 📁 项目结构

```
clawops/
├── backend/              # Go 后端
│   ├── cmd/server/       # 主入口
│   ├── internal/         # 内部包
│   │   ├── agent/        # AI Agent 核心
│   │   ├── handlers/     # HTTP Handlers
│   │   └── models/        # 数据模型
│   └── pkg/              # 公共包
│       └── database/     # 数据库
├── frontend/             # React 前端
│   └── src/
│       ├── pages/        # 页面
│       └── services/     # API 服务
├── harness/              # AI 能力配置
│   ├── skills/           # Skills 技能
│   │   ├── devops/       # 运维技能
│   │   └── software-development/  # 开发技能
│   ├── context/          # 上下文配置
│   └── constraints/      # 约束规则
└── docs/                # 文档
```

## 🔌 Skills

ClawOps 使用 Skills 扩展能力，每个 Skill 负责特定任务。

### 内置 Skills

| Skill | 说明 | 状态 |
|-------|------|------|
| deploy | 自动化部署 | ✅ |
| monitor | 监控告警 | ✅ |
| backup | 备份恢复 | ✅ |
| rollback | 回滚 | 🚧 开发中 |
| log | 日志查询 | 🚧 开发中 |

### 安装社区 Skill

```bash
# 即将支持
clawops skill install heidsoft/nginx-deploy
```

### 开发自己的 Skill

```bash
# 1. 在 harness/skills/contrib/ 创建
mkdir -p harness/skills/contrib/my-skill

# 2. 编写 SKILL.md
cat > harness/skills/contrib/my-skill/SKILL.md << 'EOF'
name: my-skill
description: 我的自定义技能
version: 1.0.0

# 技能逻辑
## When to Use
...

## Steps
1. ...
EOF

# 3. 提交 PR
```

## 🏗️ 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go + Gin + GORM |
| 前端 | React + Vite + Axios |
| 数据库 | PostgreSQL |
| AI | OpenAI / Ollama |
| 云 | 阿里云 ECS SDK |

## 📖 文档

- [快速开始](./QUICK_START.md)
- [开发指南](./DEVELOPMENT_GUIDE.md)
- [产品路线图](./ROADMAP.md)
- [CLAUDE.md](./CLAUDE.md) - AI 开发助手指南

## 🤝 贡献

欢迎贡献代码！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing`)
5. 创建 Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) 文件

## 🙏 致谢

- [OpenClaw](https://github.com/openclaw/openclaw) - AI Agent 框架
- [Harness](https://harness.io/) - 灵感来源
- [LangChain](https://github.com/langchain/langchain) - LLM 应用框架

---

**如果这个项目对你有帮助，请点个 ⭐ Star！**

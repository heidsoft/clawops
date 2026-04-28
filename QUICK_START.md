# ClawOps 快速开始

## 环境要求

- **Go**: 1.21+
- **Node.js**: 18+
- **PostgreSQL**: 14+
- **可选**: Ollama (本地 LLM) 或 OpenAI API

---

## 1. 克隆项目

```bash
git clone https://github.com/heidsoft/clawops.git
cd clawops
```

## 2. 启动 PostgreSQL

```bash
# 使用 Docker
docker run -d \
  --name clawops-db \
  -e POSTGRES_DB=clawops \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your-password \
  -p 5432:5432 \
  postgres:15
```

## 3. 配置后端

```bash
cd backend

# 复制配置模板
cp config/config.example.yaml config/config.yaml

# 编辑配置
vim config/config.yaml
```

关键配置项：

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your-password  # 修改这里
  dbname: clawops

llm:
  provider: openai        # 或 ollama
  api_key: your-api-key
  base_url: https://api.openai.com/v1

cloud:
  aliyun:
    access_key: your-access-key
    secret_key: your-secret-key
    region: cn-beijing
```

## 4. 启动后端

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端启动后会：
- 自动创建数据库表
- 监听 `http://localhost:8080`

## 5. 启动前端

新开一个终端窗口：

```bash
cd frontend
npm install
npm run dev
```

前端启动后会：
- 访问 `http://localhost:3000`
- 自动代理 API 到后端 `:8080`

## 6. 开始使用

1. 打开浏览器访问 http://localhost:3000
2. 注册/登录账号
3. 在 AI 对话框输入指令：

### 示例指令

**部署服务：**
```
部署一个 nginx 到阿里云
```

**查看部署：**
```
列出我所有的部署
```

**监控告警：**
```
查看服务器状态
```

**故障排查：**
```
网站打不开了，帮我看看什么问题
```

---

## 本地 LLM 模式 (可选)

如果你想完全离线使用：

```bash
# 1. 安装 Ollama
brew install ollama  # macOS
# 或: curl -fsSL https://ollama.com/install.sh | sh  # Linux

# 2. 下载模型
ollama pull llama3:8b

# 3. 启动 Ollama
ollama serve

# 4. 修改 backend/config/config.yaml
llm:
  provider: ollama
  base_url: http://localhost:11434
  model: llama3:8b
```

---

## 常见问题

### Q: 数据库连接失败？

```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 如果没运行，启动它
docker start clawops-db
```

### Q: 前端无法连接后端？

检查后端是否运行在 `:8080`，并确认 `frontend/vite.config.js` 的代理配置正确。

### Q: LLM API 报错？

- OpenAI: 确认 `api_key` 正确
- Ollama: 确认 `ollama serve` 已启动

---

## 下一步

- 📖 查看 [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) 了解开发规范
- 🗺️ 查看 [ROADMAP.md](ROADMAP.md) 了解产品路线图
- 🔌 查看 [harness/skills/](harness/skills/) 学习 Skill 开发

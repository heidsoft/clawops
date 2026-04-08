# ClawOps 开发计划

## 当前状态
- ✅ GitHub 仓库已创建
- ✅ 基础代码结构已提交
- ✅ README + 开发指南

## Sprint 1: 基础框架完善
**目标**: 完成可运行的基础版本

### 任务清单

| 任务 | 类型 | 优先级 | 状态 |
|------|------|--------|------|
| 数据库集成 | backend | P0 | 待开发 |
| JWT 认证 | backend | P0 | 待开发 |
| 用户管理 API | backend | P0 | 待开发 |
| 登录页面 | frontend | P0 | 待开发 |
| 部署列表页面 | frontend | P0 | 待开发 |

### 详细任务

#### Backend
1. **数据库配置** - 连接 PostgreSQL，配置 GORM
2. **用户模型** - 创建 User 模型和 CRUD
3. **认证中间件** - JWT 生成和验证
4. **部署模型** - 完善 Deployment 模型
5. **API Handler** - 实现 CRUD API

#### Frontend
1. **登录页** - 用户名/密码登录
2. **Dashboard** - 显示部署概览
3. **部署列表** - 分页展示
4. **创建部署表单** - 填写实例配置
5. **部署详情** - 显示状态和操作

## Sprint 2: 云集成
**目标**: 集成阿里云 ECS

| 任务 | 类型 | 优先级 | 状态 |
|------|--------|--------|------|
| ECS SDK 集成 | backend | P0 | 待开发 |
| 创建实例 | backend | P0 | 待开发 |
| 启动/停止实例 | backend | P0 | 待开发 |
| 前端操作按钮 | frontend | P0 | 待开发 |

## Sprint 3: 监控与日志
**目标**: 完善运维功能

| 任务 | 类型 | 优先级 | 状态 |
|------|--------|--------|------|
| 监控数据获取 | backend | P1 | 待开发 |
| 监控图表展示 | frontend | P1 | 待开发 |
| 日志收集 | backend | P1 | 待开发 |
| 日志查看器 | frontend | P1 | 待开发 |

## Sprint 4: 高级功能
**目标**: 提升用户体验

| 任务 | 类型 | 优先级 | 状态 |
|------|--------|--------|------|
| 域名自动配置 | backend | P1 | 待开发 |
| SSL 证书 | backend | P1 | 待开发 |
| 告警通知 | backend | P2 | 待开发 |
| 多云支持 | architecture | P2 | 待设计 |

---

## 开发建议

### 本地开发
1. 后端: `cd backend && go run cmd/server/main.go`
2. 前端: `cd frontend && npm run dev`
3. 数据库: 使用 Docker 启动 PostgreSQL

### 代码规范
- Backend: 遵循 Go 代码规范
- Frontend: React Hooks + 函数式组件
- 提交信息: 使用 Conventional Commits

### 测试
- Backend: `go test ./...`
- Frontend: `npm test`

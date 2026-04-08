# ClawOps 开发指南

## 项目概述
OpenClaw 部署管理系统 - 帮助用户自动化部署应用到阿里云 ECS

## 技术栈
- **后端**: Go + Gin + GORM + PostgreSQL
- **前端**: React 18 + Vite + Axios
- **云**: 阿里云 ECS SDK

## 快速开始

```bash
# 后端
cd backend && go run cmd/server/main.go

# 前端
cd frontend && npm install && npm run dev
```

## 待开发功能

### P0 - 核心功能
- [ ] 真实数据库集成 (PostgreSQL)
- [ ] JWT 用户认证
- [ ] 阿里云 ECS 真实 API 调用
- [ ] 部署实例 CRUD

### P1 - 增强功能
- [ ] 域名自动配置
- [ ] SSL 证书申请 (Let's Encrypt)
- [ ] 监控数据获取与展示
- [ ] 日志收集与查看

### P2 - 高级功能
- [ ] 自动化部署流程 (CI/CD)
- [ ] 告警通知 (邮件/短信)
- [ ] 多云支持

## API 文档
见 README.md

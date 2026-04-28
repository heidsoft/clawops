---
name: deploy
description: 自动化部署技能 - 支持 Docker、K8s、VM 部署
version: 1.0.0
author: ClawOps
metadata:
  tags: [devops, deployment, docker, kubernetes]
  related_skills: [monitor, backup, rollback]
---

# Deploy Skill

## Overview

自动化部署服务，支持多种部署方式。

## When to Use

- 需要部署新服务
- 更新现有服务版本
- 回滚到上一个版本
- 多环境部署 (dev/staging/prod)

## Supported Platforms

### Docker Compose
```bash
docker-compose up -d --build
```

### Kubernetes
```bash
kubectl apply -f deployment.yaml
kubectl set image deployment/app app=image:tag
```

### VM (SSH)
```bash
scp build.tar.gz server:/opt/app/
ssh server "cd /opt/app && ./deploy.sh"
```

## Options

| 选项 | 说明 | 默认值 |
|------|------|--------|
| platform | 部署平台 | docker |
| env | 环境 | production |
| tag | 镜像版本 | latest |
| rollback | 是否回滚 | false |

## Example

```
部署服务到生产环境
platform: kubernetes
env: production
tag: v1.2.3
```
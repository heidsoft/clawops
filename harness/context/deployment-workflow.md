# Deployment Workflow

> Context Engineering - 部署工作流规范

## 标准部署流程

```
┌─────────────────────────────────────────────────────────────────┐
│                     Deployment Workflow                         │
├─────────────────────────────────────────────────────────────────┤
│  1. Validate    →  2. Prepare  →  3. Create  →  4. Initialize  │
│                                                                  │
│  5. Deploy      →  6. Verify   →  7. Notify   →  8. Log        │
└─────────────────────────────────────────────────────────────────┘
```

## 步骤详解

### Step 1: Validate (验证)

- 验证用户输入参数
- 检查资源配额
- 确认目标地域可用性

**必填参数**:
- `user_id`: 用户 ID
- `plan`: 部署方案 (basic/standard/pro)
- `image`: Docker 镜像地址

### Step 2: Prepare (准备)

- 生成实例名称: `openclaw-{user_id}-{plan}-{timestamp}`
- 准备安全组规则
- 准备 VPC/交换机

### Step 3: Create (创建)

- 调用阿里云 CreateInstance API
- 创建 ECS 实例
- 分配公网 IP

### Step 4: Initialize (初始化)

- 安装 Docker
- 配置 Nginx
- 设置日志目录

### Step 5: Deploy (部署)

- 拉取镜像
- 启动容器
- 配置健康检查

### Step 6: Verify (验证)

- 检查端口可达性
- 验证服务响应
- 确认健康状态

### Step 7: Notify (通知)

- 发送部署成功通知
- 返回访问信息

### Step 8: Log (记录)

- 记录部署日志
- 更新数据库状态

## 部署状态

| 状态 | 说明 |
|------|------|
| `pending` | 待处理 |
| `creating` | 创建中 |
| `initializing` | 初始化中 |
| `deploying` | 部署中 |
| `running` | 运行中 |
| `failed` | 失败 |
| `stopped` | 已停止 |

## 回滚策略

如果部署失败:
1. 保留失败实例用于调试
2. 自动清理已创建的部分资源
3. 记录错误详情到数据库

## 超时设置

| 阶段 | 超时时间 |
|------|----------|
| Create | 5 分钟 |
| Initialize | 10 分钟 |
| Deploy | 15 分钟 |
| Verify | 2 分钟 |
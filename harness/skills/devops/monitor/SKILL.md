---
name: monitor
description: 监控告警技能 - 查看服务状态、告警历史、指标图表
version: 1.0.0
author: ClawOps
metadata:
  tags: [devops, monitoring, alerts, prometheus]
  related_skills: [deploy, incident, backup]
---

# Monitor Skill

## Overview

监控和告警管理，支持 Prometheus、Grafana、自定义监控。

## When to Use

- 查看服务健康状态
- 检查告警信息
- 查看指标图表
- 配置告警规则

## Commands

### 查看服务状态
```
查看所有服务状态
```

### 查看告警
```
查看最近告警
告警级别: critical
```

### 查看指标
```
CPU 使用率
内存使用情况
请求延迟 p99
```

### 告警配置
```
添加告警规则: 错误率 > 5% 持续 5 分钟
```

## Supported Tools

- Prometheus
- Grafana
- Prometheus Alertmanager
- 自定义Exporter
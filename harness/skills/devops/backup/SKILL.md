---
name: backup
description: 备份恢复技能 - 数据库、文件、配置备份
version: 1.0.0
author: ClawOps
metadata:
  tags: [devops, backup, restore, database]
  related_skills: [deploy, monitor]
---

# Backup Skill

## Overview

自动化备份和恢复，支持数据库、文件、配置等。

## When to Use

- 定期备份数据
- 故障恢复
- 迁移数据
- 查看备份历史

## Backup Types

### Database Backup
```bash
# MySQL
mysqldump -u root -p dbname > backup.sql

# PostgreSQL
pg_dump -U postgres dbname > backup.sql

# MongoDB
mongodump --db dbname --out backup/
```

### File Backup
```bash
tar -czf backup.tar.gz /data/app/
```

### Config Backup
```bash
# Kubernetes ConfigMap
kubectl get configmap -o yaml > configmap.yaml
```

## Commands

### 备份
```
备份数据库 mydb
备份 /data 目录
```

### 恢复
```
从 backup_20260101.sql 恢复
```

### 列表
```
查看所有备份
查看最近 7 天备份
```

## Retention

- 每日备份保留 7 天
- 每周备份保留 4 周
- 每月备份保留 12 月
# ClawOps Harness - 智能体工程规范

> 基于 Martin Fowler Harness Engineering 理念构建
> 参考: https://martinfowler.com/articles/exploring-gen-ai/harness-engineering-memo.html

## 核心理念

Human "On the Loop" - 人类构建和优化 harness，Agent 执行。

```
┌─────────────────────────────────────────────────────────────┐
│                      ClawOps Harness                        │
├─────────────────────────────────────────────────────────────┤
│  Context Engineering      │  知识库 + 动态上下文              │
│  Architectural Constraints │  架构约束 + Linter              │
│  Garbage Collection       │  定期清理 + 一致性检查           │
└─────────────────────────────────────────────────────────────┘
```

## 目录结构

```
harness/
├── context/           # Context Engineering
│   ├── system-prompt.md
│   ├── deployment-workflow.md
│   └── best-practices.md
├── constraints/       # Architectural Constraints
│   ├── .golangci.yml
│   ├── rules.md
│   └── arch-test/
└── gc/               # Garbage Collection
    ├── check-unused.sh
    └── verify-docs.sh
```

## 开发流程

### 1. Context Engineering

- **System Prompt**: 定义 AI 助手的角色和能力
- **Deployment Workflow**: 标准化部署流程
- **Best Practices**: 部署最佳实践文档

### 2. Architectural Constraints

- **Go Linter**: 强制代码规范
- **结构测试**: 验证模块边界
- **命名规范**: 统一命名风格

### 3. Garbage Collection

- 定期检查未使用的配置
- 验证文档一致性
- 清理过期资源

## Agent 工作流

```
1. User Request → 
2. Context Load (加载知识库) → 
3. Constraint Check (验证约束) → 
4. Execute (执行) → 
5. Validate (验证结果) → 
6. Feedback (反馈改进)
```

## 持续改进

当 Agent 遇到困难时：
- 识别缺失的工具/护栏/文档
- 补充到 harness 中
- 让 Agent 自己修复

---
Generated: 2026-04-21
Based on: Harness Engineering - first thoughts (Martin Fowler)
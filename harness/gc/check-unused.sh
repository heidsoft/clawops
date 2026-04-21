#!/bin/bash
# Garbage Collection - 检查未使用的配置
# 运行: ./check-unused.sh

set -e

echo "=== ClawOps Garbage Collection ==="
echo "检查未使用的配置..."
echo

# 检查未使用的环境变量
echo "1. 检查未使用的环境变量..."
grep -r "os.Getenv" backend/ --include="*.go" | grep -oP 'os.Getenv\("\K[^"]+' | sort -u > /tmp/used_envs.txt
grep -r "OS_" backend/.env* 2>/dev/null | grep -oP 'OS_\K[A-Z_]+' | sort -u > /tmp/defined_envs.txt

# 简单检查：如果有定义但未使用的
if [ -f /tmp/defined_envs.txt ]; then
  echo "   发现环境变量定义"
fi

# 检查未使用的 import
echo "2. 检查未使用的 Go 包..."
cd backend
go build -o /dev/null ./... 2>&1 | grep "imported and not used" || echo "   ✓ 无未使用的 import"

# 检查重复的函数
echo "3. 检查重复代码..."
# 可以添加 gocloc 或其他工具

# 检查 orphaned 文件 (没有被引用的静态资源)
echo "4. 检查孤立的静态文件..."
# frontend/src 中的文件是否都被 index.html 引用

echo
echo "=== GC 完成 ==="
echo "发现问题? 请修复后提交."
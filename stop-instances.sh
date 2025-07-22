#!/bin/bash

echo "停止所有实例..."

# 停止第一个实例
echo "停止第一个实例..."
docker-compose -p tsdd-instance1 -f docker-compose.yaml down

# 停止第二个实例
echo "停止第二个实例..."
docker-compose -p tsdd-instance2 -f docker-compose-instance2.yaml down

echo "所有实例已停止！"

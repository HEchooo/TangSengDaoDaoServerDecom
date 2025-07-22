#!/bin/bash

# 启动第一个实例
echo "启动第一个实例..."
docker-compose -p tsdd-instance1 -f docker-compose.yaml up -d

# 启动第二个实例
echo "启动第二个实例..."
docker-compose -p tsdd-instance2 -f docker-compose-instance2.yaml up -d

echo "两个实例启动完成！"
echo ""
echo "第一个实例端口映射："
echo "- WuKongIM API: 5001"
echo "- WuKongIM TCP: 5100"
echo "- WuKongIM WebSocket: 5200"
echo "- TangSeng API: 8090"
echo "- Web UI: 82"
echo "- Manager UI: 83"
echo "- Minio: 9000, 9001"
echo ""
echo "第二个实例端口映射："
echo "- WuKongIM API: 5011"
echo "- WuKongIM TCP: 5110"
echo "- WuKongIM WebSocket: 5210"
echo "- TangSeng API: 8091"
echo "- Web UI: 84"
echo "- Manager UI: 85"
echo "- Minio: 9002, 9003"

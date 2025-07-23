# 唐僧叨叨 - 快速开始指南

## 简介

唐僧叨叨是一个开源的即时通讯软件，本指南将帮助您快速部署和运行该系统。

## 快速部署

### 测试环境 (推荐新手)

1. **克隆代码并进入目录**
   ```bash
   git clone <repository-url>
   cd TangSengDaoDaoServerDecom
   ```

2. **配置环境变量**
   ```bash
   # 复制示例配置文件
   cp .env.example .env
   
   # 编辑配置文件，至少修改 EXTERNAL_IP 为您的服务器IP
   vim .env
   ```

3. **一键部署**
   ```bash
   # 构建并部署测试环境
   ./deploy.sh test -b
   
   # 或者使用 make 命令
   make deploy-test-build
   ```

4. **访问服务**
   - 业务API: http://your-ip:8090
   - 后台管理: http://your-ip:83 (用户名: superAdmin)
   - Web客户端: http://your-ip:82
   - MinIO控制台: http://your-ip:9001

### 生产环境

1. **准备AWS资源**
   - RDS MySQL实例
   - ElastiCache Redis集群
   - ECR镜像仓库访问权限

2. **配置环境变量**
   ```bash
   cp .env.prod.example .env
   vim .env  # 配置生产环境参数
   ```

3. **部署生产环境**
   ```bash
   ./deploy.sh prod -b
   ```

## 常用命令

| 操作 | 测试环境 | 生产环境 |
|------|----------|----------|
| 部署 | `make deploy-test` | `make deploy-prod` |
| 构建+部署 | `make deploy-test-build` | `make deploy-prod-build` |
| 查看日志 | `make logs-test` | `make logs-prod` |
| 查看状态 | `make status-test` | `make status-prod` |
| 停止服务 | `make stop-test` | `make stop-prod` |
| 备份数据 | `make backup` | - |
| 清理资源 | `make clean` | `make clean` |

## 服务端口说明

| 服务 | 端口 | 说明 |
|------|------|------|
| 业务API | 8090 | 主要业务接口 |
| 后台管理 | 83 | 管理后台 |
| Web客户端 | 82 | 网页版聊天 (仅测试环境) |
| WuKongIM TCP | 5100 | 长连接端口 |
| WuKongIM WebSocket | 5200 | WebSocket端口 |
| WuKongIM API | 5001 | 内部API |
| WuKongIM监控 | 5300 | 监控端口 |
| MinIO | 9000/9001 | 文件服务 (仅测试环境) |

## 默认账号

**管理员账号:**
- 用户名: `superAdmin`
- 密码: 查看 `.env` 文件中的 `TS_ADMINPWD` 配置

**测试用户 (如果有):**
- 手机号: `15900000002`
- 密码: `a1234567`

## 故障排查

### 常见问题

1. **服务无法启动**
   ```bash
   # 查看服务状态
   ./deploy.sh test -s
   
   # 查看详细日志
   ./deploy.sh test -l
   ```

2. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tlnp | grep :8090
   
   # 修改 docker-compose.yaml 中的端口映射
   ```

3. **数据库连接失败**
   ```bash
   # 检查数据库服务状态
   docker-compose ps mysql
   
   # 查看数据库日志
   docker-compose logs mysql
   ```

4. **权限问题**
   ```bash
   # 确保目录权限正确
   sudo chown -R $USER:$USER wukongim tsdd
   chmod -R 755 wukongim tsdd
   ```

### 健康检查

```bash
# 检查API是否正常
curl http://your-ip:8090/v1/ping

# 检查WuKongIM是否正常
curl http://your-ip:5001/health

# 查看服务状态
docker-compose ps
```

## 升级指南

### 测试环境升级
```bash
# 停止服务
./deploy.sh test -d

# 备份数据
./deploy.sh test --backup

# 拉取最新代码
git pull

# 重新构建并部署
./deploy.sh test -b
```

### 生产环境升级
```bash
# 备份生产数据 (根据实际情况)
# 构建新镜像
make deploy

# 停止服务
./deploy.sh prod -d

# 启动新版本
./deploy.sh prod
```

## 监控建议

- 使用 `docker stats` 监控资源使用情况
- 配置日志轮转，避免日志文件过大
- 定期备份重要数据
- 监控磁盘空间使用情况

## 安全建议

- 修改默认密码
- 配置防火墙，限制不必要的端口访问
- 使用HTTPS/WSS (生产环境)
- 定期更新系统和Docker镜像
- 配置日志审计

## 获取帮助

- 查看详细部署文档: `DEPLOYMENT.md`
- 使用部署脚本帮助: `./deploy.sh --help`
- 查看项目README: `README.md`

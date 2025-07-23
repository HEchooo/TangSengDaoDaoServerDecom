# 唐僧叨叨 Docker 部署指南

## 项目概述

唐僧叨叨是一个基于Go语言开发的开源即时通讯软件，采用微服务架构设计：

- **通讯层 (WuKongIM)**: 负责长连接维护、消息投递、消息存储等
- **业务层 (TangSengDaoDao)**: 负责用户管理、好友关系、群组管理等业务逻辑
- **Web管理端**: 后台管理系统
- **Web客户端**: 网页版聊天界面

## 服务架构

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Web Client    │  │  Web Manager    │  │   Mobile App    │
│    (Port 82)    │  │   (Port 83)     │  │                 │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                     │                     │
         └─────────────────────┼─────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│              TangSengDaoDao Business Layer                  │
│                      (Port 8090)                           │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                    WuKongIM Layer                           │
│  TCP:5100  WebSocket:5200  HTTP API:5001  Monitor:5300    │
└─────────────────────────────────────────────────────────────┘
         │                     │                     │
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│     Redis       │  │     MySQL       │  │     MinIO       │
│                 │  │                 │  │                 │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

## 部署环境要求

- Docker >= 20.0
- Docker Compose >= 2.0
- 内存: 至少 4GB RAM
- 磁盘: 至少 20GB 可用空间

## 一、测试环境部署

测试环境使用容器化的MySQL、Redis和MinIO服务，适合开发和测试使用。

### 1.1 准备环境文件

创建 `.env` 文件：

```bash
# 基础配置
EXTERNAL_IP=你的服务器IP地址

# MySQL配置
MYSQL_ROOT_PASSWORD=YourStrongPassword123
MYSQL_DATABASE=im_test

# MinIO配置
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=admin123456

# 应用配置
TS_MODE=debug
TS_FILESERVICE=minio
TS_SMSCODE=123456
TS_ADMINPWD=admin123456
TS_APPNAME=唐僧叨叨
TS_WELCOMEMESSAGE=欢迎使用唐僧叨叨

# WuKongIM配置
WK_MODE=debug
WK_CONVERSATION_ON=true
WK_DATASOURCE_CHANNELINFOON=true
WK_TOKENAUTHON=true
WK_WHITELISTOFFOFPERSON=false
```

### 1.2 创建测试环境 docker-compose.yml

```yaml
version: '3.1'
services:
  # Redis服务
  redis:
    image: redis:7.2.3
    restart: always
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 1s
      timeout: 3s
      retries: 30
    volumes:
      - redis_data:/data

  # MySQL数据库
  mysql:
    image: mysql:8.0.33
    command: --default-authentication-plugin=mysql_native_password
    restart: always
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 30s
      timeout: 10s
      retries: 5
    environment:
      - TZ=Asia/Shanghai
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_DATABASE=${MYSQL_DATABASE}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"

  # MinIO文件服务
  minio:
    image: minio/minio:RELEASE.2023-07-18T17-49-40Z
    restart: always
    command: "server /data --console-address ':9001'"
    environment:
      - MINIO_ROOT_USER=${MINIO_ROOT_USER}
      - MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD}
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3
    volumes:
      - minio_data:/data
    ports:
      - "9000:9000"
      - "9001:9001"

  # WuKongIM通讯服务
  wukongim:
    image: registry.cn-shanghai.aliyuncs.com/wukongim/wukongim:v1.2
    restart: always
    depends_on:
      redis:
        condition: service_healthy
      mysql:
        condition: service_healthy
    ports:
      - "5001:5001"  # HTTP API端口
      - "5100:5100"  # TCP长连接端口
      - "5200:5200"  # WebSocket端口
      - "5300:5300"  # 监控端口
    volumes:
      - ./wukongim:/root/wukongim
    env_file:
      - .env
    environment:
      - WK_EXTERNAL_IP=${EXTERNAL_IP}
      - WK_WEBHOOK_GRPCADDR=tangsengdaodaoserverdecom:6979
      - WK_DATASOURCE_ADDR=http://tangsengdaodaoserverdecom:8090/v1/datasource

  # 唐僧叨叨业务服务
  tangsengdaodaoserverdecom:
    image: tangsengdaodaoserverdecom:latest
    restart: always
    command: "api"
    depends_on:
      - redis
      - mysql
      - wukongim
      - minio
    ports:
      - "8090:8090"
    volumes:
      - ./tsdd:/home/tsdddata
      - ./tsdd/configs:/home/configs
    env_file:
      - .env
    environment:
      - TS_DB_MYSQLADDR=root:${MYSQL_ROOT_PASSWORD}@tcp(mysql:3306)/${MYSQL_DATABASE}?charset=utf8mb4&parseTime=true&loc=Local
      - TS_DB_REDISADDR=redis:6379
      - TS_WUKONGIM_APIURL=http://wukongim:5001
      - TS_MINIO_URL=http://minio:9000
      - TS_EXTERNAL_IP=${EXTERNAL_IP}
    healthcheck:
      test: "wget -q -Y off -O /dev/null http://localhost:8090/v1/ping > /dev/null 2>&1"
      interval: 10s
      timeout: 10s
      retries: 3

  # Web管理端
  tangsengdaodaomanager:
    image: registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaomanager:latest
    restart: always
    depends_on:
      - tangsengdaodaoserverdecom
    environment:
      - API_URL=http://tangsengdaodaoserverdecom:8090/
    ports:
      - "83:80"

  # Web客户端
  tangsengdaodaoweb:
    image: registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaoweb:latest
    restart: always
    depends_on:
      - tangsengdaodaoserverdecom
    environment:
      - API_URL=http://tangsengdaodaoserverdecom:8090/
    ports:
      - "82:80"

  # 数据库管理工具 (仅测试环境)
  adminer:
    image: adminer:latest
    restart: always
    ports:
      - "8306:8080"

volumes:
  redis_data:
  mysql_data:
  minio_data:
```

### 1.3 部署步骤

```bash
# 1. 构建项目镜像
make build

# 2. 创建必要的目录
mkdir -p wukongim tsdd/configs tsdd/logs

# 3. 启动服务
docker-compose up -d

# 4. 查看服务状态
docker-compose ps

# 5. 查看日志
docker-compose logs -f tangsengdaodaoserverdecom
```

## 二、生产环境部署

生产环境使用AWS托管的RDS MySQL和ElastiCache Redis服务，具有更好的可靠性和性能。

### 2.1 准备环境文件

创建生产环境 `.env` 文件：

```bash
# 基础配置
EXTERNAL_IP=你的生产服务器IP

# AWS RDS MySQL配置
MYSQL_HOST=your-rds-endpoint.region.rds.amazonaws.com:3306
MYSQL_PASSWORD=YourProductionPassword
MYSQL_DATABASE=im_production

# AWS ElastiCache Redis配置
REDIS_HOST=your-redis-cluster.cache.amazonaws.com:6379

# MinIO/S3配置 (可选择使用AWS S3)
TS_FILESERVICE=s3
# 或者使用自建MinIO
# TS_FILESERVICE=minio
MINIO_ROOT_USER=your-access-key
MINIO_ROOT_PASSWORD=your-secret-key

# 生产环境配置
TS_MODE=release
TS_ADMINPWD=YourSecureAdminPassword123
TS_APPNAME=YourAppName
TS_EXTERNAL_BASEURL=https://your-domain.com
WK_EXTERNAL_WSSADDR=wss://your-ws-domain.com

# WuKongIM生产配置
WK_MODE=release
WK_CONVERSATION_ON=true
WK_DATASOURCE_CHANNELINFOON=true
WK_TOKENAUTHON=true

# 推送服务配置
ECHOOO_PUSH_SERVERADDRESSES=your-push-servers
TS_PUSH_FIREBASE_PROJECTID=your-firebase-project
TS_PUSH_FIREBASE_JSONPATH=/home/configs/firebase-key.json
```

### 2.2 创建生产环境 docker-compose.prod.yml

```yaml
version: '3.1'
services:
  # WuKongIM通讯服务
  wukongim:
    image: registry.cn-shanghai.aliyuncs.com/wukongim/wukongim:v1.2
    restart: always
    ports:
      - "5001:5001"  # HTTP API端口 (内网)
      - "5100:5100"  # TCP长连接端口 (外网)
      - "5200:5200"  # WebSocket端口 (外网)
      - "5300:5300"  # 监控端口
    volumes:
      - ./wukongim:/root/wukongim
    env_file:
      - .env
    environment:
      - WK_EXTERNAL_IP=${EXTERNAL_IP}
      - WK_WEBHOOK_GRPCADDR=tangsengdaodaoserverdecom:6979
      - WK_DATASOURCE_ADDR=http://tangsengdaodaoserverdecom:8090/v1/datasource
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  # 唐僧叨叨业务服务
  tangsengdaodaoserverdecom:
    image: 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
    restart: always
    command: "api"
    depends_on:
      - wukongim
    ports:
      - "8090:8090"
    volumes:
      - ./tsdd:/home/tsdddata
      - ./tsdd/configs:/home/configs
    env_file:
      - .env
    environment:
      - TS_DB_MYSQLADDR=im:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST})/${MYSQL_DATABASE}?charset=utf8mb4&parseTime=true&loc=Local
      - TS_DB_REDISADDR=${REDIS_HOST}
      - TS_WUKONGIM_APIURL=http://wukongim:5001
      - TS_EXTERNAL_IP=${EXTERNAL_IP}
      - TS_FILESERVICE=${TS_FILESERVICE}
      - TS_MINIO_ACCESSKEYID=${MINIO_ROOT_USER}
      - TS_MINIO_SECRETACCESSKEY=${MINIO_ROOT_PASSWORD}
      - TS_ECHOOO_PUSH_SERVERADDRESSES=${ECHOOO_PUSH_SERVERADDRESSES}
    healthcheck:
      test: "wget -q -Y off -O /dev/null http://localhost:8090/v1/ping > /dev/null 2>&1"
      interval: 10s
      timeout: 10s
      retries: 3
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  # Web管理端
  tangsengdaodaomanager:
    image: registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaomanager:latest
    restart: always
    depends_on:
      - tangsengdaodaoserverdecom
    environment:
      - API_URL=http://tangsengdaodaoserverdecom:8090/
    ports:
      - "83:80"
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "2"
```

### 2.3 生产环境部署步骤

```bash
# 1. 确保AWS服务已准备就绪
# - RDS MySQL实例已创建并可访问
# - ElastiCache Redis集群已创建并可访问
# - ECR镜像仓库访问权限已配置

# 2. 登录AWS ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 851725583589.dkr.ecr.us-east-1.amazonaws.com

# 3. 构建并推送镜像
make deploy

# 4. 创建生产环境目录结构
mkdir -p wukongim tsdd/configs tsdd/logs

# 5. 上传配置文件 (如果有自定义配置)
cp your-firebase-key.json tsdd/configs/

# 6. 启动生产服务
docker-compose -f docker-compose.prod.yml up -d

# 7. 验证服务状态
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f
```

## 三、配置说明

### 3.1 核心环境变量说明

| 变量名 | 说明 | 测试环境示例 | 生产环境示例 |
|--------|------|-------------|-------------|
| `EXTERNAL_IP` | 服务器外网IP | `192.168.1.100` | `3.216.154.243` |
| `TS_MODE` | 运行模式 | `debug` | `release` |
| `TS_DB_MYSQLADDR` | MySQL连接地址 | `root:pass@tcp(mysql:3306)/db` | `user:pass@tcp(rds-host:3306)/db` |
| `TS_DB_REDISADDR` | Redis连接地址 | `redis:6379` | `elasticache-host:6379` |
| `TS_FILESERVICE` | 文件服务类型 | `minio` | `s3` 或 `minio` |

### 3.2 文件服务配置

**MinIO配置**:
```bash
TS_FILESERVICE=minio
TS_MINIO_URL=http://minio:9000  # 或外部MinIO地址
TS_MINIO_ACCESSKEYID=your-access-key
TS_MINIO_SECRETACCESSKEY=your-secret-key
```

**AWS S3配置**:
```bash
TS_FILESERVICE=s3
# S3相关配置通过AWS凭证或IAM角色配置
```

### 3.3 推送服务配置

```bash
# Firebase推送
TS_PUSH_FIREBASE_PROJECTID=your-project-id
TS_PUSH_FIREBASE_JSONPATH=/home/configs/firebase-key.json

# iOS推送
TS_PUSH_APNS_TOPIC=com.yourcompany.app
TS_PUSH_APNS_CERT=/home/configs/push.p12
TS_PUSH_APNS_PASSWORD=cert-password
```

## 四、运维管理

### 4.1 服务监控

```bash
# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]

# 查看资源使用情况
docker stats

# 重启特定服务
docker-compose restart tangsengdaodaoserverdecom
```

### 4.2 数据备份

```bash
# MySQL数据备份 (测试环境)
docker-compose exec mysql mysqldump -u root -p${MYSQL_ROOT_PASSWORD} ${MYSQL_DATABASE} > backup.sql

# Redis数据备份 (测试环境)
docker-compose exec redis redis-cli BGSAVE

# MinIO数据备份
docker-compose exec minio mc mirror /data /backup/minio-data
```

### 4.3 常用维护命令

```bash
# 更新服务镜像
docker-compose pull
docker-compose up -d

# 清理未使用的镜像
docker image prune -f

# 查看磁盘使用情况
df -h
docker system df
```

## 五、故障排查

### 5.1 常见问题

**服务无法启动**:
```bash
# 检查端口占用
netstat -tlnp | grep :8090

# 检查配置文件
docker-compose config

# 查看详细错误日志
docker-compose logs --tail=100 tangsengdaodaoserverdecom
```

**数据库连接失败**:
```bash
# 测试数据库连接
docker-compose exec tangsengdaodaoserverdecom ping mysql

# 检查数据库状态
docker-compose exec mysql mysqladmin -u root -p status
```

**文件上传失败**:
```bash
# 检查MinIO服务状态
curl http://localhost:9000/minio/health/live

# 检查MinIO访问权限
docker-compose exec minio mc admin info local
```

### 5.2 性能优化

**数据库优化**:
- 配置适当的连接池大小
- 启用慢查询日志分析
- 定期执行数据库维护操作

**Redis优化**:
- 配置适当的内存使用策略
- 启用持久化配置
- 监控内存使用情况

**应用优化**:
- 调整日志级别 (生产环境使用info或warn)
- 配置适当的并发处理参数
- 启用应用性能监控

## 六、安全建议

1. **网络安全**:
   - 限制数据库端口仅内网访问
   - 使用防火墙限制不必要的端口暴露
   - 配置HTTPS/WSS加密传输

2. **认证安全**:
   - 使用强密码策略
   - 定期更换数据库密码
   - 配置API访问权限控制

3. **数据安全**:
   - 定期备份重要数据
   - 启用数据库加密
   - 配置日志审计

通过本部署指南，您应该能够成功部署唐僧叨叨即时通讯系统。如遇到问题，请参考故障排查章节或查看相关日志进行诊断。

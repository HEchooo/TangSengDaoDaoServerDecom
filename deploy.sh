#!/bin/bash

# 唐僧叨叨部署脚本
# 支持测试环境和生产环境部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 显示使用说明
show_usage() {
    cat << EOF
唐僧叨叨部署脚本

用法: $0 [选项] [环境]

环境:
  test        测试环境部署 (使用容器化数据库)
  prod        生产环境部署 (使用AWS RDS/ElastiCache)

选项:
  -h, --help     显示此帮助信息
  -b, --build    构建镜像
  -d, --down     停止并删除服务
  -l, --logs     查看服务日志
  -s, --status   查看服务状态
  --backup       备份数据 (仅测试环境)
  --clean        清理未使用的Docker资源

示例:
  $0 test                # 部署测试环境
  $0 prod -b             # 构建镜像并部署生产环境
  $0 test -l             # 查看测试环境日志
  $0 test --backup       # 备份测试环境数据

EOF
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    log_info "依赖检查完成"
}

# 检查环境文件
check_env_file() {
    local env_type=$1
    local env_file=".env"
    local example_file=""
    
    if [[ "$env_type" == "test" ]]; then
        example_file=".env.example"
    else
        example_file=".env.prod.example"
    fi
    
    if [[ ! -f "$env_file" ]]; then
        log_warn "环境文件 $env_file 不存在"
        if [[ -f "$example_file" ]]; then
            log_info "复制示例环境文件..."
            cp "$example_file" "$env_file"
            log_warn "请编辑 $env_file 文件，配置正确的环境变量"
            exit 1
        else
            log_error "示例环境文件 $example_file 不存在"
            exit 1
        fi
    fi
    
    log_info "环境文件检查完成"
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    mkdir -p wukongim/data
    mkdir -p wukongim/logs
    mkdir -p tsdd/configs
    mkdir -p tsdd/logs
    
    # 设置目录权限
    chmod 755 wukongim tsdd
    chmod 755 wukongim/data wukongim/logs
    chmod 755 tsdd/configs tsdd/logs
    
    log_info "目录创建完成"
}

# 构建镜像
build_image() {
    log_info "构建应用镜像..."
    
    if ! make build; then
        log_error "镜像构建失败"
        exit 1
    fi
    
    log_info "镜像构建完成"
}

# 部署测试环境
deploy_test() {
    local compose_file="docker-compose.yaml"
    
    log_info "部署测试环境..."
    
    check_env_file "test"
    create_directories
    
    # 启动服务
    log_info "启动服务..."
    docker-compose -f "$compose_file" up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 30
    
    # 检查服务状态
    check_services "$compose_file"
    
    log_info "测试环境部署完成!"
    show_access_info "test"
}

# 部署生产环境
deploy_prod() {
    local compose_file="docker-compose.prod.yml"
    
    log_info "部署生产环境..."
    
    check_env_file "prod"
    create_directories
    
    # 检查AWS ECR登录
    log_info "检查AWS ECR访问..."
    if ! aws sts get-caller-identity &> /dev/null; then
        log_warn "AWS CLI未配置或无权限，尝试使用ECR登录..."
        if ! aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 851725583589.dkr.ecr.us-east-1.amazonaws.com; then
            log_error "AWS ECR登录失败，请检查AWS凭证配置"
            exit 1
        fi
    fi
    
    # 拉取最新镜像
    log_info "拉取最新镜像..."
    docker-compose -f "$compose_file" pull
    
    # 启动服务
    log_info "启动服务..."
    docker-compose -f "$compose_file" up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 30
    
    # 检查服务状态
    check_services "$compose_file"
    
    log_info "生产环境部署完成!"
    show_access_info "prod"
}

# 检查服务状态
check_services() {
    local compose_file=$1
    
    log_info "检查服务状态..."
    docker-compose -f "$compose_file" ps
    
    # 检查健康状态
    local unhealthy_services=$(docker-compose -f "$compose_file" ps --format "table {{.Name}}\t{{.Status}}" | grep -v "Up" | grep -v "Name" || true)
    
    if [[ -n "$unhealthy_services" ]]; then
        log_warn "以下服务状态异常:"
        echo "$unhealthy_services"
        log_info "查看服务日志以获取更多信息: $0 $ENV_TYPE -l"
    else
        log_info "所有服务运行正常"
    fi
}

# 显示访问信息
show_access_info() {
    local env_type=$1
    local external_ip=$(grep EXTERNAL_IP .env | cut -d'=' -f2)
    
    echo
    log_info "=== 服务访问信息 ==="
    echo "业务API: http://${external_ip}:8090"
    echo "后台管理: http://${external_ip}:83"
    echo "WebSocket: ws://${external_ip}:5200"
    echo "WuKongIM API: http://${external_ip}:5001"
    echo "WuKongIM监控: http://${external_ip}:5300"
    
    if [[ "$env_type" == "test" ]]; then
        echo "Web客户端: http://${external_ip}:82"
        echo "MinIO控制台: http://${external_ip}:9001"
        echo "数据库管理: http://${external_ip}:8306"
    fi
    
    echo
    log_info "默认管理员账号:"
    echo "用户名: superAdmin"
    echo "密码: 查看 .env 文件中的 TS_ADMINPWD 配置"
    echo
}

# 查看日志
show_logs() {
    local compose_file=$1
    local service=${2:-""}
    
    if [[ -n "$service" ]]; then
        log_info "显示 $service 服务日志..."
        docker-compose -f "$compose_file" logs -f "$service"
    else
        log_info "显示所有服务日志..."
        docker-compose -f "$compose_file" logs -f
    fi
}

# 停止服务
stop_services() {
    local compose_file=$1
    
    log_info "停止服务..."
    docker-compose -f "$compose_file" down
    
    log_info "服务已停止"
}

# 备份数据
backup_data() {
    local backup_dir="backup/$(date +%Y%m%d_%H%M%S)"
    
    log_info "备份数据到 $backup_dir..."
    mkdir -p "$backup_dir"
    
    # 备份MySQL数据
    if docker-compose ps mysql | grep -q "Up"; then
        log_info "备份MySQL数据..."
        docker-compose exec -T mysql mysqldump -u root -p"${MYSQL_ROOT_PASSWORD}" "${MYSQL_DATABASE}" > "$backup_dir/mysql_backup.sql"
    fi
    
    # 备份Redis数据
    if docker-compose ps redis | grep -q "Up"; then
        log_info "备份Redis数据..."
        docker-compose exec -T redis redis-cli --rdb - > "$backup_dir/redis_backup.rdb"
    fi
    
    # 备份MinIO数据
    if [[ -d "miniodata" ]]; then
        log_info "备份MinIO数据..."
        tar -czf "$backup_dir/minio_backup.tar.gz" miniodata/
    fi
    
    # 备份配置文件
    log_info "备份配置文件..."
    cp -r tsdd/configs "$backup_dir/" 2>/dev/null || true
    cp .env "$backup_dir/" 2>/dev/null || true
    
    log_info "数据备份完成: $backup_dir"
}

# 清理Docker资源
clean_docker() {
    log_info "清理未使用的Docker资源..."
    
    docker image prune -f
    docker container prune -f
    docker network prune -f
    docker volume prune -f
    
    log_info "清理完成"
}

# 主函数
main() {
    local env_type=""
    local action=""
    local build_flag=false
    local service=""
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            test|prod)
                env_type="$1"
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -b|--build)
                build_flag=true
                shift
                ;;
            -d|--down)
                action="down"
                shift
                ;;
            -l|--logs)
                action="logs"
                shift
                ;;
            -s|--status)
                action="status"
                shift
                ;;
            --backup)
                action="backup"
                shift
                ;;
            --clean)
                action="clean"
                shift
                ;;
            *)
                service="$1"
                shift
                ;;
        esac
    done
    
    # 检查环境类型
    if [[ -z "$env_type" ]] && [[ "$action" != "clean" ]]; then
        log_error "请指定环境类型: test 或 prod"
        show_usage
        exit 1
    fi
    
    # 检查依赖
    check_dependencies
    
    # 设置compose文件
    local compose_file="docker-compose.yaml"
    if [[ "$env_type" == "prod" ]]; then
        compose_file="docker-compose.prod.yml"
    fi
    
    # 构建镜像
    if [[ "$build_flag" == true ]]; then
        build_image
    fi
    
    # 执行操作
    case $action in
        "down")
            stop_services "$compose_file"
            ;;
        "logs")
            show_logs "$compose_file" "$service"
            ;;
        "status")
            check_services "$compose_file"
            ;;
        "backup")
            if [[ "$env_type" != "test" ]]; then
                log_error "备份功能仅支持测试环境"
                exit 1
            fi
            backup_data
            ;;
        "clean")
            clean_docker
            ;;
        *)
            # 默认部署操作
            if [[ "$env_type" == "test" ]]; then
                deploy_test
            elif [[ "$env_type" == "prod" ]]; then
                deploy_prod
            fi
            ;;
    esac
}

# 执行主函数
main "$@"

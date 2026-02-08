#!/bin/bash

# ╔════════════════════════════════════════════════════════════╗
# ║              一键启动脚本 - NodeLoc 发卡系统                 ║
# ╚════════════════════════════════════════════════════════════╝

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 欢迎信息
clear
cat << "EOF"
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║        NodeLoc 社区发卡系统 - 一键部署脚本                   ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
EOF
echo ""

# 检查 Docker 和 Docker Compose
print_info "检查系统依赖..."

if ! command -v docker &> /dev/null; then
    print_error "Docker 未安装，请先安装 Docker"
    echo "安装指南: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose 未安装，请先安装 Docker Compose"
    echo "安装指南: https://docs.docker.com/compose/install/"
    exit 1
fi

print_success "Docker 和 Docker Compose 已安装"
echo ""

# 检查 .env 文件
if [ ! -f ".env" ]; then
    print_warning ".env 文件不存在"
    
    if [ -f ".env.production" ]; then
        print_info "复制 .env.production 为 .env"
        cp .env.production .env
        print_success ".env 文件已创建"
    else
        print_error ".env.production 文件也不存在，请手动创建配置文件"
        exit 1
    fi
    
    echo ""
    print_warning "=========================================="
    print_warning "重要: 请编辑 .env 文件配置必要参数"
    print_warning "=========================================="
    echo ""
    echo "必须配置的参数:"
    echo "  - DB_ROOT_PASSWORD (MySQL root 密码)"
    echo "  - DB_PASSWORD (应用数据库密码)"
    echo "  - SESSION_SECRET (32位随机字符串)"
    echo "  - NODELOC_CLIENT_ID"
    echo "  - NODELOC_CLIENT_SECRET"
    echo "  - NODELOC_REDIRECT_URI"
    echo ""
    read -p "按回车键继续编辑配置文件..." 
    
    ${EDITOR:-nano} .env
fi

# 加载环境变量
source .env

# 检查必要配置
print_info "检查配置文件..."
REQUIRED_VARS=(
    "DB_PASSWORD"
    "NODELOC_CLIENT_ID"
    "NODELOC_CLIENT_SECRET"
    "NODELOC_REDIRECT_URI"
)

MISSING_VARS=()
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        MISSING_VARS+=("$var")
    fi
done

if [ ${#MISSING_VARS[@]} -gt 0 ]; then
    print_error "以下必要配置项未设置:"
    for var in "${MISSING_VARS[@]}"; do
        echo "  - $var"
    done
    echo ""
    print_info "请编辑 .env 文件并设置这些配置项"
    exit 1
fi

print_success "配置检查通过"
echo ""

# 选择部署方式
print_info "选择部署方式:"
echo "  1) 从源码构建（需要时间，适合开发）"
echo "  2) 使用 Docker Hub 镜像（快速，适合生产）"
echo ""
read -p "请选择 [1/2] (默认: 2): " DEPLOY_METHOD
DEPLOY_METHOD=${DEPLOY_METHOD:-2}

if [ "$DEPLOY_METHOD" = "1" ]; then
    COMPOSE_FILE="docker-compose.yml"
    print_info "将从源码构建镜像..."
else
    COMPOSE_FILE="docker-compose.prod.yml"
    print_info "将使用 Docker Hub 预构建镜像..."
fi

echo ""

# 停止现有服务
if docker-compose ps | grep -q "Up"; then
    print_warning "检测到正在运行的服务"
    read -p "是否停止现有服务? [y/N]: " STOP_EXISTING
    if [[ $STOP_EXISTING =~ ^[Yy]$ ]]; then
        print_info "停止现有服务..."
        docker-compose down
        print_success "服务已停止"
    fi
    echo ""
fi

# 拉取/构建镜像
print_info "=========================================="
print_info "准备镜像..."
print_info "=========================================="

if [ "$DEPLOY_METHOD" = "1" ]; then
    print_info "开始构建镜像（可能需要几分钟）..."
    docker-compose -f "$COMPOSE_FILE" build
else
    print_info "拉取镜像..."
    docker-compose -f "$COMPOSE_FILE" pull
fi

print_success "镜像准备完成"
echo ""

# 启动服务
print_info "=========================================="
print_info "启动服务..."
print_info "=========================================="

docker-compose -f "$COMPOSE_FILE" up -d

# 等待服务启动
print_info "等待服务启动..."
sleep 5

# 检查服务状态
print_info "检查服务状态..."
docker-compose ps

echo ""
print_info "=========================================="
print_success "部署完成！"
print_info "=========================================="
echo ""

# 获取访问地址
PORT=${PORT:-3000}
echo -e "${GREEN}访问地址:${NC}"
echo "  http://localhost:${PORT}"
echo "  http://$(hostname -I | awk '{print $1}'):${PORT}"
echo ""

echo -e "${YELLOW}默认管理员账号:${NC}"
echo "  用户名: admin"
echo "  密码: admin123"
echo "  ${RED}⚠️  请立即登录并修改密码！${NC}"
echo ""

echo -e "${BLUE}常用命令:${NC}"
echo "  查看日志:    docker-compose logs -f"
echo "  重启服务:    docker-compose restart"
echo "  停止服务:    docker-compose down"
echo "  查看状态:    docker-compose ps"
echo ""

print_info "查看详细部署文档: ./docker-deploy.md"
echo ""

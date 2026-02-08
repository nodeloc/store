#!/bin/bash

# ╔════════════════════════════════════════════════════════════╗
# ║        发布镜像到 Docker Hub (nodeloc/faka)                 ║
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

# 显示标题
clear
cat << "EOF"
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║     发布 NodeLoc 发卡系统镜像到 Docker Hub                   ║
║               组织: nodeloc                                ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
EOF
echo ""

# 配置
DOCKER_ORG="nodeloc"
VERSION="${1:-latest}"
BACKEND_IMAGE="${DOCKER_ORG}/faka-backend:${VERSION}"
FRONTEND_IMAGE="${DOCKER_ORG}/faka-frontend:${VERSION}"

print_info "=========================================="
print_info "发布配置"
print_info "=========================================="
echo "  Docker Hub 组织: ${DOCKER_ORG}"
echo "  版本标签: ${VERSION}"
echo "  后端镜像: ${BACKEND_IMAGE}"
echo "  前端镜像: ${FRONTEND_IMAGE}"
print_info "=========================================="
echo ""

# 检查 Docker 登录状态
print_info "检查 Docker Hub 登录状态..."
if ! docker info | grep -q "Username"; then
    print_warning "未登录到 Docker Hub"
    print_info "正在登录..."
    docker login
    echo ""
fi
print_success "Docker Hub 已登录"
echo ""

# 确认发布
read -p "确认发布以上配置的镜像? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "取消发布"
    exit 0
fi
echo ""

# 构建后端镜像
print_info "=========================================="
print_info "步骤 1/4: 构建后端镜像"
print_info "=========================================="
docker build \
    --platform linux/amd64,linux/arm64 \
    -t "${BACKEND_IMAGE}" \
    -f Dockerfile \
    . || {
    print_error "后端镜像构建失败"
    exit 1
}

# 如果不是 latest，同时打 latest 标签
if [ "$VERSION" != "latest" ]; then
    docker tag "${BACKEND_IMAGE}" "${DOCKER_ORG}/faka-backend:latest"
fi

print_success "后端镜像构建完成"
echo ""

# 构建前端镜像
print_info "=========================================="
print_info "步骤 2/4: 构建前端镜像"
print_info "=========================================="
print_info "提示: 如果构建慢，已配置使用淘宝 npm 镜像源"
docker build \
    -t "${FRONTEND_IMAGE}" \
    -f frontend/Dockerfile \
    frontend/ || {
    print_error "前端镜像构建失败"
    print_info "如果遇到 npm install 卡住，请检查网络或使用代理"
    exit 1
}

# 如果不是 latest，同时打 latest 标签
if [ "$VERSION" != "latest" ]; then
    docker tag "${FRONTEND_IMAGE}" "${DOCKER_ORG}/faka-frontend:latest"
fi

print_success "前端镜像构建完成"
echo ""

# 推送后端镜像
print_info "=========================================="
print_info "步骤 3/4: 推送后端镜像"
print_info "=========================================="
docker push "${BACKEND_IMAGE}" || {
    print_error "后端镜像推送失败"
    exit 1
}

if [ "$VERSION" != "latest" ]; then
    docker push "${DOCKER_ORG}/faka-backend:latest"
fi

print_success "后端镜像推送完成"
echo ""

# 推送前端镜像
print_info "=========================================="
print_info "步骤 4/4: 推送前端镜像"
print_info "=========================================="
docker push "${FRONTEND_IMAGE}" || {
    print_error "前端镜像推送失败"
    exit 1
}

if [ "$VERSION" != "latest" ]; then
    docker push "${DOCKER_ORG}/faka-frontend:latest"
fi

print_success "前端镜像推送完成"
echo ""

# 显示发布信息
print_info "=========================================="
print_success "✅ 所有镜像已成功发布！"
print_info "=========================================="
echo ""
echo "已发布的镜像："
echo "  📦 ${BACKEND_IMAGE}"
echo "  📦 ${FRONTEND_IMAGE}"

if [ "$VERSION" != "latest" ]; then
    echo "  📦 ${DOCKER_ORG}/faka-backend:latest"
    echo "  📦 ${DOCKER_ORG}/faka-frontend:latest"
fi

echo ""
echo "用户可以通过以下命令使用："
echo ""
echo "# 拉取镜像"
echo "docker pull ${BACKEND_IMAGE}"
echo "docker pull ${FRONTEND_IMAGE}"
echo ""
echo "# 使用 Docker Compose 部署"
echo "docker-compose -f docker-compose.prod.yml up -d"
echo ""

print_info "Docker Hub 地址："
echo "  https://hub.docker.com/r/${DOCKER_ORG}/faka-backend"
echo "  https://hub.docker.com/r/${DOCKER_ORG}/faka-frontend"
echo ""

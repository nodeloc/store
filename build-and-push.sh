#!/bin/bash

# ╔════════════════════════════════════════════════════════════╗
# ║        Docker 镜像构建和推送脚本                             ║
# ╚════════════════════════════════════════════════════════════╝

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
DOCKER_USERNAME="${DOCKER_USERNAME:-nodeloc}"
VERSION="${VERSION:-latest}"
BACKEND_IMAGE="${BACKEND_IMAGE:-faka-backend}"
FRONTEND_IMAGE="${FRONTEND_IMAGE:-faka-frontend}"

# 镜像完整名称
BACKEND_FULL_NAME="${DOCKER_USERNAME}/${BACKEND_IMAGE}:${VERSION}"
FRONTEND_FULL_NAME="${DOCKER_USERNAME}/${FRONTEND_IMAGE}:${VERSION}"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示使用说明
show_usage() {
    cat << EOF
使用方法:
    $0 [选项]

选项:
    -u, --username USERNAME     Docker Hub 用户名 (默认: nodeloc)
    -v, --version VERSION       镜像版本标签 (默认: latest)
    -b, --backend-only          仅构建后端镜像
    -f, --frontend-only         仅构建前端镜像
    -n, --no-push              仅构建不推送
    -h, --help                  显示此帮助信息

示例:
    # 构建并推送所有镜像
    $0

    # 指定用户名和版本
    $0 -u myusername -v v1.0.0

    # 仅构建后端镜像
    $0 -b

    # 构建但不推送
    $0 -n

环境变量:
    DOCKER_USERNAME            Docker Hub 用户名
    VERSION                    镜像版本标签
EOF
}

# 解析命令行参数
BUILD_BACKEND=true
BUILD_FRONTEND=true
PUSH_IMAGES=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--username)
            DOCKER_USERNAME="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -b|--backend-only)
            BUILD_FRONTEND=false
            shift
            ;;
        -f|--frontend-only)
            BUILD_BACKEND=false
            shift
            ;;
        -n|--no-push)
            PUSH_IMAGES=false
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "未知选项: $1"
            show_usage
            exit 1
            ;;
    esac
done

# 显示配置信息
print_info "==========================================<br/>"
print_info "Docker 镜像构建配置"
print_info "==========================================<br/>"
print_info "Docker Hub 用户名: ${DOCKER_USERNAME}"
print_info "版本标签: ${VERSION}"
print_info "构建后端: ${BUILD_BACKEND}"
print_info "构建前端: ${BUILD_FRONTEND}"
print_info "推送镜像: ${PUSH_IMAGES}"
print_info "==========================================<br/>"
echo ""

# 检查 Docker 是否已登录
check_docker_login() {
    if [ "$PUSH_IMAGES" = true ]; then
        print_info "检查 Docker 登录状态..."
        if ! docker info | grep -q "Username"; then
            print_warning "未登录到 Docker Hub"
            print_info "请先登录: docker login"
            read -p "是否现在登录? (y/n) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                docker login
            else
                print_error "取消操作"
                exit 1
            fi
        else
            print_success "Docker 已登录"
        fi
    fi
}

# 构建后端镜像
build_backend() {
    print_info "==========================================<br/>"
    print_info "开始构建后端镜像..."
    print_info "==========================================<br/>"
    
    docker build \
        --platform linux/amd64,linux/arm64 \
        -t "${BACKEND_FULL_NAME}" \
        -f Dockerfile \
        .
    
    # 同时打 latest 标签
    if [ "$VERSION" != "latest" ]; then
        docker tag "${BACKEND_FULL_NAME}" "${DOCKER_USERNAME}/${BACKEND_IMAGE}:latest"
    fi
    
    print_success "后端镜像构建完成: ${BACKEND_FULL_NAME}"
}

# 构建前端镜像
build_frontend() {
    print_info "=========================================="
    print_info "开始构建前端镜像..."
    print_info "=========================================="
    
    docker build \
        --platform linux/amd64,linux/arm64 \
        -t "${FRONTEND_FULL_NAME}" \
        -f frontend/Dockerfile \
        frontend/
    
    # 同时打 latest 标签
    if [ "$VERSION" != "latest" ]; then
        docker tag "${FRONTEND_FULL_NAME}" "${DOCKER_USERNAME}/${FRONTEND_IMAGE}:latest"
    fi
    
    print_success "前端镜像构建完成: ${FRONTEND_FULL_NAME}"
}

# 推送镜像
push_images() {
    print_info "=========================================="
    print_info "开始推送镜像到 Docker Hub..."
    print_info "=========================================="
    
    if [ "$BUILD_BACKEND" = true ]; then
        print_info "推送后端镜像..."
        docker push "${BACKEND_FULL_NAME}"
        if [ "$VERSION" != "latest" ]; then
            docker push "${DOCKER_USERNAME}/${BACKEND_IMAGE}:latest"
        fi
        print_success "后端镜像推送完成"
    fi
    
    if [ "$BUILD_FRONTEND" = true ]; then
        print_info "推送前端镜像..."
        docker push "${FRONTEND_FULL_NAME}"
        if [ "$VERSION" != "latest" ]; then
            docker push "${DOCKER_USERNAME}/${FRONTEND_IMAGE}:latest"
        fi
        print_success "前端镜像推送完成"
    fi
}

# 显示镜像信息
show_images_info() {
    print_info "=========================================="
    print_info "构建的镜像信息"
    print_info "=========================================="
    
    if [ "$BUILD_BACKEND" = true ]; then
        echo -e "${GREEN}后端镜像:${NC}"
        echo "  - ${BACKEND_FULL_NAME}"
        if [ "$VERSION" != "latest" ]; then
            echo "  - ${DOCKER_USERNAME}/${BACKEND_IMAGE}:latest"
        fi
        docker images | grep "${DOCKER_USERNAME}/${BACKEND_IMAGE}" | head -2
    fi
    
    echo ""
    
    if [ "$BUILD_FRONTEND" = true ]; then
        echo -e "${GREEN}前端镜像:${NC}"
        echo "  - ${FRONTEND_FULL_NAME}"
        if [ "$VERSION" != "latest" ]; then
            echo "  - ${DOCKER_USERNAME}/${FRONTEND_IMAGE}:latest"
        fi
        docker images | grep "${DOCKER_USERNAME}/${FRONTEND_IMAGE}" | head -2
    fi
    
    echo ""
}

# 主流程
main() {
    # 检查登录状态
    check_docker_login
    
    # 构建镜像
    if [ "$BUILD_BACKEND" = true ]; then
        build_backend
    fi
    
    if [ "$BUILD_FRONTEND" = true ]; then
        build_frontend
    fi
    
    # 推送镜像
    if [ "$PUSH_IMAGES" = true ]; then
        push_images
    else
        print_warning "跳过推送镜像（使用 -n 选项）"
    fi
    
    # 显示镜像信息
    show_images_info
    
    # 显示使用说明
    print_info "=========================================="
    print_success "所有操作完成！"
    print_info "=========================================="
    echo ""
    print_info "使用镜像："
    
    if [ "$BUILD_BACKEND" = true ]; then
        echo "  docker pull ${BACKEND_FULL_NAME}"
    fi
    
    if [ "$BUILD_FRONTEND" = true ]; then
        echo "  docker pull ${FRONTEND_FULL_NAME}"
    fi
    
    echo ""
    print_info "或使用 Docker Compose 部署："
    echo "  docker-compose -f docker-compose.prod.yml up -d"
    echo ""
}

# 错误处理
trap 'print_error "发生错误，脚本终止"; exit 1' ERR

# 运行主流程
main

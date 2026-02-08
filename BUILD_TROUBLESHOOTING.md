# 🚀 Docker 构建优化指南

## 常见问题与解决方案

### 1. npm install 一直卡住 ⚠️

**问题原因：** npm 默认使用国外镜像源，在国内访问很慢

**解决方案：**

#### 方案 A：使用优化后的 Dockerfile（推荐）✅

已在 `frontend/Dockerfile` 中配置淘宝镜像源：

```dockerfile
RUN npm config set registry https://registry.npmmirror.com
RUN npm install --legacy-peer-deps --no-audit --progress=false
```

重新构建即可：

```bash
./publish.sh
```

#### 方案 B：使用 Docker 构建参数

如果还是慢，可以在构建时使用代理：

```bash
# 使用 HTTP 代理
docker build \
  --build-arg HTTP_PROXY=http://proxy.example.com:8080 \
  --build-arg HTTPS_PROXY=http://proxy.example.com:8080 \
  -t nodeloc/faka-frontend:latest \
  ./frontend

# 使用 socks5 代理
docker build \
  --build-arg HTTP_PROXY=socks5://127.0.0.1:1080 \
  --build-arg HTTPS_PROXY=socks5://127.0.0.1:1080 \
  -t nodeloc/faka-frontend:latest \
  ./frontend
```

#### 方案 C：预先安装依赖

在本地先安装依赖，然后复制到 Docker：

```bash
cd frontend
npm install --registry=https://registry.npmmirror.com
cd ..

# 修改 Dockerfile 复制 node_modules
```

#### 方案 D：使用 pnpm（更快）

修改 `frontend/Dockerfile`：

```dockerfile
FROM node:18-alpine as builder

RUN npm install -g pnpm
RUN pnpm config set registry https://registry.npmmirror.com

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build
```

---

### 2. Docker 构建很慢

**优化技巧：**

#### 使用 BuildKit

```bash
# 启用 BuildKit
export DOCKER_BUILDKIT=1

# 构建时显示详细输出
docker build --progress=plain -t nodeloc/faka-backend:latest .
```

#### 使用多阶段构建缓存

```bash
# 构建时保留缓存
docker build --cache-from nodeloc/faka-backend:latest -t nodeloc/faka-backend:latest .
```

#### 清理 Docker 缓存

如果构建出错，清理缓存重试：

```bash
# 清理构建缓存
docker builder prune -af

# 重新构建
./publish.sh
```

---

### 3. 后端构建优化

#### 使用 Go 模块代理（国内）

修改 `Dockerfile`：

```dockerfile
# 在构建阶段添加
ENV GOPROXY=https://goproxy.cn,direct
ENV GOPRIVATE=github.com/your-private-repo

RUN go mod download
```

#### 使用缓存挂载

```dockerfile
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /app/faka .
```

---

### 4. 镜像体积优化

#### 当前镜像大小

```bash
docker images | grep nodeloc/faka
```

#### 优化建议

1. **使用 Alpine 基础镜像** ✅ 已使用
2. **多阶段构建** ✅ 已使用
3. **清理不必要的文件**

```dockerfile
# 在构建阶段
RUN npm install && npm cache clean --force
RUN go build -ldflags="-s -w" -o app
```

4. **使用 .dockerignore** ✅ 已创建

---

### 5. 网络问题诊断

#### 测试网络连接

```bash
# 进入构建容器测试
docker run --rm -it node:18-alpine sh

# 测试 npm 源
npm config get registry
npm config set registry https://registry.npmmirror.com

# 测试下载速度
npm install lodash --verbose
```

#### 使用国内 Docker Hub 镜像

配置 Docker daemon（`/etc/docker/daemon.json`）：

```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com"
  ]
}
```

重启 Docker：

```bash
sudo systemctl restart docker
```

---

## 快速解决方案

### 最快的方式：使用已构建的镜像

如果构建一直失败，可以直接使用已发布的镜像：

```bash
# 不构建，直接拉取
docker pull nodeloc/faka-backend:latest
docker pull nodeloc/faka-frontend:latest

# 使用 docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

---

## 构建时间参考

正常情况下的构建时间：

- **后端构建**: 2-5 分钟
- **前端构建**: 3-8 分钟（取决于网络）
- **总时间**: 5-15 分钟

如果超过 30 分钟，说明有问题，按上述方案排查。

---

## 推荐的完整构建命令

```bash
# 1. 启用 BuildKit
export DOCKER_BUILDKIT=1

# 2. 清理缓存（如果之前构建失败）
docker builder prune -af

# 3. 构建后端（使用国内 Go 代理）
docker build \
  --build-arg GOPROXY=https://goproxy.cn,direct \
  --progress=plain \
  -t nodeloc/faka-backend:latest \
  .

# 4. 构建前端（已配置淘宝源）
docker build \
  --progress=plain \
  -t nodeloc/faka-frontend:latest \
  ./frontend

# 5. 推送镜像
docker push nodeloc/faka-backend:latest
docker push nodeloc/faka-frontend:latest
```

---

## 在国内服务器上的最佳实践

```bash
# 1. 配置系统级代理（如果有）
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080

# 2. 配置 npm 镜像
npm config set registry https://registry.npmmirror.com

# 3. 配置 Go 代理
go env -w GOPROXY=https://goproxy.cn,direct

# 4. 配置 Docker 镜像加速
# 编辑 /etc/docker/daemon.json

# 5. 然后构建
./publish.sh
```

---

## 故障排查命令

```bash
# 查看构建日志
docker build --progress=plain --no-cache -t test ./frontend 2>&1 | tee build.log

# 进入构建阶段调试
docker run --rm -it node:18-alpine sh

# 手动测试构建步骤
cd /tmp
git clone https://github.com/nodeloc/store.git
cd store/frontend
npm config set registry https://registry.npmmirror.com
npm install --verbose

# 查看网络连接
ping registry.npmjs.org
ping registry.npmmirror.com
curl -I https://registry.npmmirror.com
```

---

## 联系支持

如果以上方案都无法解决问题：

1. 查看完整的构建日志
2. 提交 Issue 并附带日志
3. 在 NodeLoc 论坛求助
4. 考虑使用已构建的镜像

GitHub Issues: https://github.com/nodeloc/store/issues

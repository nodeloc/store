# 📦 发布镜像到 Docker Hub (nodeloc/faka)

## 🎯 镜像信息

### Docker Hub 组织
- **组织名**: `nodeloc`
- **项目名**: `faka`

### 镜像列表
- **后端镜像**: `nodeloc/faka-backend`
- **前端镜像**: `nodeloc/faka-frontend`

### Docker Hub 地址
- 后端: https://hub.docker.com/r/nodeloc/faka-backend
- 前端: https://hub.docker.com/r/nodeloc/faka-frontend

---

## 🚀 快速发布

### 方法 1: 使用自动化脚本（推荐）

```bash
# 发布 latest 版本
./publish.sh

# 发布指定版本
./publish.sh v1.0.0
```

脚本会自动：
1. 检查 Docker 登录状态
2. 构建后端镜像
3. 构建前端镜像
4. 推送所有镜像到 Docker Hub
5. 如果是版本标签，同时推送 latest

### 方法 2: 手动发布

```bash
# 1. 登录 Docker Hub
docker login

# 2. 构建镜像
docker build -t nodeloc/faka-backend:latest .
docker build -t nodeloc/faka-frontend:latest ./frontend

# 3. 推送镜像
docker push nodeloc/faka-backend:latest
docker push nodeloc/faka-frontend:latest
```

### 方法 3: 使用通用构建脚本

```bash
# 使用 build-and-push.sh
./build-and-push.sh -u nodeloc -v latest
```

---

## 📋 发布前检查清单

### ✅ 必须完成的步骤

- [ ] 代码已提交到 Git
- [ ] 已测试所有功能正常
- [ ] 已更新 README.md
- [ ] 已登录 Docker Hub (`docker login`)
- [ ] 确认有 nodeloc 组织的推送权限

### ✅ 推荐完成的步骤

- [ ] 更新版本号
- [ ] 更新 CHANGELOG（如果有）
- [ ] 测试 Docker 镜像构建成功
- [ ] 本地测试镜像运行正常

---

## 🔢 版本管理

### 版本号规范

使用语义化版本号：`vMAJOR.MINOR.PATCH`

- **MAJOR**: 重大更新，不兼容的 API 变更
- **MINOR**: 新功能，向后兼容
- **PATCH**: Bug 修复，向后兼容

### 发布不同版本

```bash
# 开发版本
./publish.sh dev

# 测试版本
./publish.sh beta

# 正式版本
./publish.sh v1.0.0

# 最新版本（默认）
./publish.sh latest
```

### 同时打多个标签

```bash
# 构建镜像
docker build -t nodeloc/faka-backend:v1.0.0 .
docker build -t nodeloc/faka-backend:v1.0 .
docker build -t nodeloc/faka-backend:v1 .
docker build -t nodeloc/faka-backend:latest .

# 推送所有标签
docker push nodeloc/faka-backend:v1.0.0
docker push nodeloc/faka-backend:v1.0
docker push nodeloc/faka-backend:v1
docker push nodeloc/faka-backend:latest
```

---

## 🏗️ 多平台构建

### 使用 Docker Buildx（支持 ARM）

```bash
# 1. 创建 builder
docker buildx create --name faka-builder --use

# 2. 启动 builder
docker buildx inspect --bootstrap

# 3. 构建并推送多平台镜像
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t nodeloc/faka-backend:latest \
  --push \
  .

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t nodeloc/faka-frontend:latest \
  --push \
  ./frontend
```

### 支持的平台

- `linux/amd64` - x86_64 架构（常见服务器）
- `linux/arm64` - ARM64 架构（树莓派、Apple Silicon）
- `linux/arm/v7` - ARM v7 架构（可选）

---

## 📝 发布流程

### 完整发布流程

1. **准备发布**
   ```bash
   # 确保代码已提交
   git status
   git add .
   git commit -m "Release v1.0.0"
   git tag v1.0.0
   ```

2. **登录 Docker Hub**
   ```bash
   docker login
   # 输入 nodeloc 组织的凭证
   ```

3. **执行发布**
   ```bash
   ./publish.sh v1.0.0
   ```

4. **验证发布**
   ```bash
   # 拉取刚发布的镜像
   docker pull nodeloc/faka-backend:v1.0.0
   docker pull nodeloc/faka-frontend:v1.0.0
   
   # 测试运行
   docker-compose -f docker-compose.prod.yml up -d
   ```

5. **推送 Git 标签**
   ```bash
   git push origin v1.0.0
   git push origin main
   ```

6. **创建 GitHub Release**
   - 访问 GitHub Release 页面
   - 创建新的 Release
   - 填写更新日志
   - 附加相关文件

---

## 🔍 验证发布

### 检查镜像是否发布成功

```bash
# 检查镜像标签
docker manifest inspect nodeloc/faka-backend:latest
docker manifest inspect nodeloc/faka-frontend:latest

# 拉取镜像测试
docker pull nodeloc/faka-backend:latest
docker pull nodeloc/faka-frontend:latest

# 检查镜像大小
docker images | grep nodeloc/faka
```

### 测试部署

```bash
# 创建测试目录
mkdir -p /tmp/faka-test
cd /tmp/faka-test

# 下载配置文件
curl -O https://raw.githubusercontent.com/nodeloc/store/main/docker-compose.prod.yml
curl -O https://raw.githubusercontent.com/nodeloc/store/main/.env.production
cp .env.production .env

# 编辑配置
nano .env

# 启动测试
docker-compose -f docker-compose.prod.yml up -d

# 查看日志
docker-compose logs -f

# 清理测试
docker-compose down -v
cd -
rm -rf /tmp/faka-test
```

---

## 📊 镜像信息

### 查看已发布的镜像

访问 Docker Hub：
- https://hub.docker.com/r/nodeloc/faka-backend/tags
- https://hub.docker.com/r/nodeloc/faka-frontend/tags

### 镜像大小预估

- **后端镜像**: ~20-30 MB（Alpine + Go 二进制）
- **前端镜像**: ~30-40 MB（Alpine + Nginx + 静态文件）

---

## 🔐 权限管理

### Docker Hub 组织权限

确保你有以下权限：
- ✅ nodeloc 组织成员
- ✅ 可以推送镜像
- ✅ 可以管理标签

### 添加团队成员

如果需要其他人发布：
1. 访问 https://hub.docker.com/orgs/nodeloc/teams
2. 邀请成员
3. 分配适当的权限

---

## 🛠️ 常见问题

### 1. 推送失败：权限被拒绝

```bash
# 确保已登录
docker login

# 检查登录状态
docker info | grep Username

# 确认有 nodeloc 组织权限
```

### 2. 构建失败

```bash
# 清理缓存重新构建
docker builder prune -af
docker build --no-cache -t nodeloc/faka-backend:latest .
```

### 3. 镜像过大

```bash
# 检查镜像层
docker history nodeloc/faka-backend:latest

# 优化 Dockerfile
# - 使用多阶段构建
# - 合并 RUN 命令
# - 清理缓存
```

### 4. 多平台构建失败

```bash
# 确保安装了 QEMU
docker run --privileged --rm tonistiigi/binfmt --install all

# 重新创建 builder
docker buildx rm faka-builder
docker buildx create --name faka-builder --use
docker buildx inspect --bootstrap
```

---

## 📞 获取帮助

### Docker Hub 支持
- 文档: https://docs.docker.com/docker-hub/
- 支持: https://hub.docker.com/support

### 项目支持
- GitHub Issues: https://github.com/nodeloc/store/issues
- NodeLoc 论坛: https://www.nodeloc.com

---

## 📅 发布历史

记录主要版本的发布历史：

| 版本 | 发布日期 | 主要更新 |
|------|----------|----------|
| v1.0.0 | 2026-01-20 | 初始版本发布 |
| latest | 持续更新 | 最新开发版本 |

---

## 🎉 发布后

### 更新文档

- [ ] 更新 README.md 中的版本号
- [ ] 更新部署文档
- [ ] 在社区发布更新公告

### 通知用户

可以通过以下方式通知用户更新：

1. **GitHub Release 通知**
2. **NodeLoc 论坛公告**
3. **项目 README 更新说明**

### 示例更新通知

```markdown
## 🎉 新版本发布: v1.0.0

### 更新内容
- 新功能 A
- 修复 Bug B
- 性能优化 C

### 更新方式
使用 Docker Compose 的用户：
\`\`\`bash
docker-compose pull
docker-compose up -d
\`\`\`

详细说明: [查看更新日志](链接)
```

---

## 📚 相关文档

- [README.md](README.md) - 项目主文档
- [docker-deploy.md](docker-deploy.md) - Docker 部署文档
- [DOCKER.md](DOCKER.md) - Docker 配置说明
- [build-and-push.sh](build-and-push.sh) - 通用构建脚本
- [publish.sh](publish.sh) - 快速发布脚本

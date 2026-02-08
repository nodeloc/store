# 🐳 Docker 一键部署配置文件说明

## 📁 相关文件

### Docker 配置文件

1. **docker-compose.yml** - 标准 Docker Compose 配置（包含 MySQL）
   - 适用于本地开发和从源码构建
   - 包含完整的 MySQL、后端、前端服务
   - 支持健康检查和依赖管理

2. **docker-compose.prod.yml** - 生产环境配置
   - 使用 Docker Hub 预构建镜像
   - 优化的生产环境配置
   - 快速部署，无需本地构建

3. **Dockerfile** - 后端镜像构建文件
   - 多阶段构建优化镜像大小
   - 基于 Alpine Linux
   - 非 root 用户运行

4. **frontend/Dockerfile** - 前端镜像构建文件
   - Node.js 构建 + Nginx 服务
   - 生产优化构建

5. **.dockerignore** / **frontend/.dockerignore** - Docker 构建忽略文件
   - 减少构建上下文大小
   - 加快构建速度

### 配置文件

1. **.env.production** - 生产环境配置示例
   - 包含所有必要的环境变量
   - 带有详细的注释说明
   - 复制为 .env 使用

2. **env.example** - 开发环境配置示例
   - 适用于本地开发
   - 使用外部 MySQL

### 脚本文件

1. **quick-start.sh** - 一键启动脚本 ⭐
   - 自动检查依赖
   - 交互式配置
   - 自动启动服务
   - **推荐新手使用**

2. **build-and-push.sh** - 构建和推送脚本
   - 构建 Docker 镜像
   - 推送到 Docker Hub
   - 支持版本标签
   - **适合维护者使用**

### 文档文件

1. **docker-deploy.md** - Docker 部署完整文档
   - 详细的部署步骤
   - 常用命令参考
   - 故障排查指南
   - 生产环境建议

2. **README.md** - 项目主文档
   - 项目介绍
   - 快速开始
   - API 文档

---

## 🚀 快速开始指南

### 方法 1：使用一键脚本（推荐新手）

```bash
git clone https://github.com/nodeloc/store.git
cd store
./quick-start.sh
```

### 方法 2：使用 Docker Hub 镜像（最快）

```bash
# 1. 准备配置文件
cp .env.production .env
nano .env  # 编辑必要配置

# 2. 启动服务
docker-compose -f docker-compose.prod.yml up -d

# 3. 查看日志
docker-compose logs -f
```

### 方法 3：从源码构建

```bash
# 1. 准备配置文件
cp .env.production .env
nano .env  # 编辑必要配置

# 2. 构建并启动
docker-compose up -d

# 3. 查看日志
docker-compose logs -f
```

---

## 🔧 配置说明

### 必须配置的环境变量

```env
# 数据库密码（必改！）
DB_ROOT_PASSWORD=your_strong_root_password
DB_PASSWORD=your_strong_db_password

# Session 密钥（32位随机字符串）
SESSION_SECRET=your_random_32_char_secret_here

# NodeLoc OAuth（从 nodeloc.com 获取）
NODELOC_CLIENT_ID=your_oauth_client_id
NODELOC_CLIENT_SECRET=your_oauth_client_secret
NODELOC_REDIRECT_URI=https://your-domain.com/auth/callback
```

### 可选配置

```env
# 前端端口（默认 3000）
PORT=3000

# NodeLoc Payment（如需支付功能）
PAYMENT_ID=pay_xxxxxxxxxx
PAYMENT_SECRET=your_payment_secret
PAYMENT_CALLBACK_URI=https://your-domain.com/payment/callback

# 自定义 Docker 镜像
DOCKER_IMAGE_BACKEND=yourusername/faka-backend:latest
DOCKER_IMAGE_FRONTEND=yourusername/faka-frontend:latest
```

---

## 📦 服务说明

### MySQL（数据库）
- **端口**: 3306（仅本地访问）
- **数据持久化**: mysql_data 卷
- **默认数据库**: faka
- **字符集**: utf8mb4

### Backend（后端 API）
- **内部端口**: 8080
- **框架**: Go + Gin
- **健康检查**: /api/health
- **数据持久化**: uploads 卷

### Frontend（前端）
- **外部端口**: 3000（可配置）
- **Web 服务器**: Nginx
- **框架**: Vue 3 + Vite

---

## 🔍 服务状态检查

```bash
# 查看所有服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f mysql

# 进入容器
docker-compose exec backend sh
docker-compose exec mysql mysql -u faka -p
```

---

## 🛠 常见操作

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart backend
```

### 更新镜像

```bash
# 拉取最新镜像
docker-compose pull

# 重新创建容器
docker-compose up -d
```

### 备份数据

```bash
# 备份数据库
docker-compose exec mysql mysqldump -u root -p${DB_ROOT_PASSWORD} faka > backup.sql

# 备份上传文件
tar -czf uploads-backup.tar.gz uploads/
```

### 清理

```bash
# 停止并删除容器
docker-compose down

# 停止并删除容器和数据卷（危险！）
docker-compose down -v

# 清理未使用的镜像
docker image prune -a
```

---

## 🌐 生产环境建议

### 1. 使用反向代理

推荐使用 Nginx 或 Caddy 作为反向代理，并配置 SSL：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. 配置防火墙

```bash
# 仅开放必要端口
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

### 3. 定期备份

设置自动备份脚本（crontab）：

```bash
0 2 * * * /path/to/backup-script.sh
```

### 4. 监控服务

使用 Docker 健康检查和日志监控：

```bash
# 查看服务健康状态
docker ps --format "table {{.Names}}\t{{.Status}}"

# 设置日志轮转
docker-compose logs --tail=1000 > logs/app.log
```

---

## 📚 更多文档

- **完整部署文档**: [docker-deploy.md](docker-deploy.md)
- **项目文档**: [README.md](README.md)
- **API 文档**: README.md 的 API 部分

---

## 💡 技巧

1. **快速查看所有配置选项**
   ```bash
   cat .env.production
   ```

2. **测试配置是否正确**
   ```bash
   docker-compose config
   ```

3. **查看容器资源使用**
   ```bash
   docker stats
   ```

4. **导出镜像到文件**
   ```bash
   docker save -o faka-images.tar nodeloc/faka-backend nodeloc/faka-frontend
   ```

5. **从文件加载镜像**
   ```bash
   docker load -i faka-images.tar
   ```

---

## 🆘 获取帮助

- **GitHub Issues**: https://github.com/nodeloc/store/issues
- **NodeLoc 论坛**: https://www.nodeloc.com
- **查看日志**: `docker-compose logs -f`

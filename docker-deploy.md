# Docker 一键部署指南

## 📦 快速部署（推荐）

### 方式一：使用 Docker Compose 直接部署

#### 1. 克隆项目

```bash
git clone https://github.com/nodeloc/store.git
cd store
```

#### 2. 配置环境变量

```bash
# 复制配置文件
cp .env.production .env

# 编辑配置文件
nano .env
```

**必须修改的配置：**
- `DB_ROOT_PASSWORD` - MySQL root 密码
- `DB_PASSWORD` - 应用数据库密码
- `SESSION_SECRET` - Session 密钥（32位随机字符串）
- `NODELOC_CLIENT_ID` - NodeLoc OAuth 客户端 ID
- `NODELOC_CLIENT_SECRET` - NodeLoc OAuth 客户端密钥
- `NODELOC_REDIRECT_URI` - OAuth 回调地址（如 https://your-domain.com/auth/callback）

**可选配置：**
- `PAYMENT_ID` 和 `PAYMENT_SECRET` - 如需支付功能
- `PORT` - 修改默认端口（默认 3000）

#### 3. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

#### 4. 访问系统

打开浏览器访问：`http://your-server-ip:3000`

默认管理员账号：
- 用户名：`admin`
- 密码：`admin123`

**⚠️ 重要：首次登录后请立即修改管理员密码！**

---

## 🐋 使用 Docker Hub 镜像部署

如果项目已经构建并推送到 Docker Hub，可以直接拉取镜像部署。

### 1. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: nodeloc-faka-mysql
    restart: unless-stopped
    environment:
      - MYSQL_ROOT_PASSWORD=${DB_ROOT_PASSWORD}
      - MYSQL_DATABASE=${DB_NAME:-faka}
      - MYSQL_USER=${DB_USER:-faka}
      - MYSQL_PASSWORD=${DB_PASSWORD}
      - TZ=Asia/Shanghai
    ports:
      - "${DB_PORT:-3306}:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    command: --default-authentication-plugin=mysql_native_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    networks:
      - faka-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${DB_ROOT_PASSWORD}"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    image: nodeloc/faka-backend:latest  # 从 Docker Hub 拉取
    container_name: nodeloc-faka-backend
    restart: unless-stopped
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=${DB_USER:-faka}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME:-faka}
      - NODELOC_URL=${NODELOC_URL:-https://www.nodeloc.com}
      - NODELOC_CLIENT_ID=${NODELOC_CLIENT_ID}
      - NODELOC_CLIENT_SECRET=${NODELOC_CLIENT_SECRET}
      - NODELOC_REDIRECT_URI=${NODELOC_REDIRECT_URI}
      - SESSION_SECRET=${SESSION_SECRET}
      - PAYMENT_ID=${PAYMENT_ID}
      - PAYMENT_SECRET=${PAYMENT_SECRET}
      - PAYMENT_CALLBACK_URI=${PAYMENT_CALLBACK_URI}
      - SERVER_PORT=8080
      - GIN_MODE=release
    volumes:
      - uploads:/app/uploads
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - faka-network

  frontend:
    image: nodeloc/faka-frontend:latest  # 从 Docker Hub 拉取
    container_name: nodeloc-faka-frontend
    restart: unless-stopped
    ports:
      - "${PORT:-3000}:80"
    depends_on:
      - backend
    networks:
      - faka-network

volumes:
  mysql_data:
  uploads:

networks:
  faka-network:
    driver: bridge
```

### 2. 创建 .env 文件

```bash
# 复制配置
cp .env.production .env

# 编辑配置
nano .env
```

### 3. 启动服务

```bash
docker-compose up -d
```

---

## 🔨 构建并推送到 Docker Hub

如果你想构建自己的镜像并推送到 Docker Hub：

### 1. 登录 Docker Hub

```bash
docker login
```

### 2. 构建镜像

```bash
# 构建后端镜像
docker build -t yourusername/faka-backend:latest .

# 构建前端镜像
docker build -t yourusername/faka-frontend:latest ./frontend
```

### 3. 推送镜像

```bash
# 推送后端镜像
docker push yourusername/faka-backend:latest

# 推送前端镜像
docker push yourusername/faka-frontend:latest
```

### 4. 修改 docker-compose.yml

将 `docker-compose.yml` 中的镜像名称修改为你的镜像：

```yaml
backend:
  image: yourusername/faka-backend:latest
  # ...

frontend:
  image: yourusername/faka-frontend:latest
  # ...
```

或者在 `.env` 文件中设置：

```env
DOCKER_IMAGE_BACKEND=yourusername/faka-backend:latest
DOCKER_IMAGE_FRONTEND=yourusername/faka-frontend:latest
```

---

## 📋 常用命令

### 服务管理

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart backend

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
```

### 数据管理

```bash
# 备份数据库
docker-compose exec mysql mysqldump -u root -p${DB_ROOT_PASSWORD} faka > backup.sql

# 恢复数据库
docker-compose exec -T mysql mysql -u root -p${DB_ROOT_PASSWORD} faka < backup.sql

# 进入 MySQL 容器
docker-compose exec mysql mysql -u root -p

# 查看数据卷
docker volume ls

# 删除数据卷（危险操作！）
docker-compose down -v
```

### 更新服务

```bash
# 拉取最新镜像
docker-compose pull

# 重新构建并启动
docker-compose up -d --build

# 仅重新构建特定服务
docker-compose build backend
docker-compose up -d backend
```

---

## 🔧 配置 Nginx 反向代理（可选）

如果你想通过域名访问，可以在宿主机配置 Nginx：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS（如果有 SSL 证书）
    # return 301 https://$server_name$request_uri;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTPS 配置
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

---

## 🐛 故障排查

### 1. 容器无法启动

```bash
# 查看详细日志
docker-compose logs

# 检查配置文件
docker-compose config
```

### 2. 数据库连接失败

```bash
# 检查 MySQL 容器状态
docker-compose ps mysql

# 查看 MySQL 日志
docker-compose logs mysql

# 测试数据库连接
docker-compose exec mysql mysql -u faka -p -e "SHOW DATABASES;"
```

### 3. 后端服务健康检查失败

```bash
# 查看后端日志
docker-compose logs backend

# 进入容器检查
docker-compose exec backend sh
wget -O- http://localhost:8080/api/health
```

### 4. 端口冲突

如果默认端口被占用，修改 `.env` 文件中的 `PORT` 变量：

```env
PORT=8080  # 修改为其他端口
```

---

## 📊 监控和维护

### 资源使用情况

```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
docker system df
```

### 清理无用资源

```bash
# 清理未使用的镜像
docker image prune -a

# 清理未使用的容器
docker container prune

# 清理未使用的数据卷
docker volume prune

# 清理所有未使用的资源
docker system prune -a
```

---

## 🔐 安全建议

1. **修改默认密码**
   - 立即修改 MySQL root 密码
   - 修改管理员账号密码

2. **使用强密码**
   - 数据库密码至少 16 位
   - Session Secret 至少 32 位随机字符串

3. **配置防火墙**
   - 仅开放必要端口（80, 443）
   - 限制 3306 端口仅本地访问

4. **定期备份**
   - 每天自动备份数据库
   - 定期备份上传的文件

5. **启用 HTTPS**
   - 使用 Let's Encrypt 免费证书
   - 配置 Nginx 反向代理

6. **监控日志**
   - 定期检查异常访问
   - 启用日志轮转避免磁盘占满

---

## 📞 获取帮助

如遇到问题，请：

1. 查看 [GitHub Issues](https://github.com/nodeloc/store/issues)
2. 访问 [NodeLoc 论坛](https://www.nodeloc.com)
3. 查看项目 [README.md](README.md)

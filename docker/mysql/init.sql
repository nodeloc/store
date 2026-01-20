-- NodeLoc 发卡系统数据库初始化脚本
-- 此脚本会在 MySQL 容器首次启动时自动执行

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 授予用户权限（如果需要从外部连接）
-- GRANT ALL PRIVILEGES ON faka.* TO 'faka'@'%';
-- FLUSH PRIVILEGES;

-- 说明：数据库表会由应用程序通过 GORM 自动创建和迁移
-- 无需在此手动创建表结构

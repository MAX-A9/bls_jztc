# Docker部署说明

本文档说明如何在Linux服务器上使用Docker部署应用。

## 目录结构

```
manifest/docker/
├── Dockerfile          # Docker镜像构建文件
├── docker-compose.yml  # Docker Compose配置文件
├── deploy.sh           # 部署脚本
├── docker.sh           # Docker构建准备脚本
├── wait-for-mysql.sh   # 等待MySQL就绪脚本
├── start-app.sh        # 应用启动脚本
├── stop.sh             # 停止服务脚本
├── restart.sh          # 重启服务脚本
└── status.sh           # 查看服务状态脚本
```

## 环境要求

- CentOS/RHEL/Ubuntu/Debian Linux
- Docker 20.10+
- Docker Compose 2.0+
- 至少2GB内存
- 至少20GB磁盘空间

## 部署步骤

1. 使用PowerShell构建应用
   ```powershell
   cd demo
   ./build.ps1  # 生成build目录
   ```

2. 将build目录上传到Linux服务器
   ```bash
   # 在Windows上使用SCP上传
   scp -r build/ user@server:/build
   
   # 或使用SFTP上传
   # 或使用其他文件传输工具
   ```

3. 在Linux服务器上进入Docker目录并部署
   ```bash
   cd /build/manifest/docker
   # 确保脚本使用LF行尾
   dos2unix *.sh
   # 添加执行权限
   chmod +x *.sh
   # 执行部署脚本
   ./deploy.sh
   ```

> 注意：部署脚本只能在Linux服务器上运行，不适用于Windows环境

## 行尾符号问题处理

如果在Windows系统上编辑脚本后上传到Linux服务器，可能会出现行尾符号不兼容的问题，导致`$'\r': command not found`等错误。处理方法：

1. 安装dos2unix工具
   ```bash
   # CentOS/RHEL
   sudo yum install dos2unix
   
   # Ubuntu/Debian
   sudo apt-get install dos2unix
   ```

2. 转换所有脚本文件
   ```bash
   dos2unix *.sh
   ```

3. 如果无法安装dos2unix，可以使用sed命令
   ```bash
   sed -i 's/\r$//' *.sh
   ```

## 部署详解

### 1. 配置和准备阶段

部署脚本在执行时会：
- 检查并安装Docker和Docker Compose
- 检查端口(8000和3306)是否被占用
- 自动查找和复制配置文件和SQL脚本
- 创建必要的目录(logs, data等)
- 运行docker.sh脚本准备Docker构建环境
  - 生成等待MySQL就绪的脚本
  - 生成应用启动脚本

### 2. 构建和启动阶段

- 应用使用多阶段构建减小镜像体积
- 数据库使用MySQL 8.0官方镜像
- 应用和数据库使用专用网络通信
- 数据库数据持久化到卷中

### 3. 启动和验证阶段

- 应用启动前会自动等待MySQL就绪
- 支持环境变量注入配置
- 健康检查确保服务正常运行

## 服务管理

### 停止服务
```bash
./stop.sh
```

### 重启服务
```bash
./restart.sh
```

### 查看服务状态
```bash
./status.sh
```

## 访问地址

- 应用访问地址：http://服务器IP:8000
- Swagger文档地址：http://服务器IP:8000/swagger

## 数据库配置

- 数据库主机：mysql
- 数据库端口：3306
- 数据库名：go-jztc
- 用户名：go-jztc
- 密码：BeAdwSKMbeiDpkHw

## 目录挂载

- 配置文件: `../config` → `/app/manifest/config`
- 日志文件: `../../logs` → `/app/logs`
- SQL脚本: `../sql` → `/docker-entrypoint-initdb.d`
- 数据库数据: `mysql-data` (Docker卷)

## 常见问题

1. **配置文件找不到**
   ```
   panic: configuration not found, did you miss the configuration file...
   ```
   解决方案:
   - 确保配置文件存在于正确位置
   - 部署脚本会自动查找并复制配置文件
   - 如果问题依然存在，可以手动创建目录并复制配置：
     ```bash
     mkdir -p /build/manifest/config
     cp /build/manifest/config/config.yaml /build/manifest/config/
     ```

2. **行尾符号问题**
   ```
   $'\r': command not found
   ```
   解决方案:
   - 使用`dos2unix`转换所有脚本文件的行尾符号
   - 或使用`sed -i 's/\r$//' *.sh`命令

3. **无法连接数据库**
   ```
   [ERROR] 无法连接到数据库，超过最大重试次数
   ```
   解决方案:
   - 检查MySQL容器是否正常运行: `docker ps | grep jztc-mysql`
   - 查看MySQL日志: `docker logs jztc-mysql`
   - 确认网络配置是否正确，特别是host设置

4. **服务无法启动**
   解决方案:
   - 查看应用日志: `docker logs jztc-app`
   - 检查配置文件是否正确
   - 使用`docker exec -it jztc-app sh`进入容器检查文件

5. **端口被占用**
   ```
   [ERROR] 端口 8000 已被占用，请先释放该端口
   ```
   解决方案:
   - 查找并关闭占用端口的进程: `lsof -i :8000`
   - 修改`docker-compose.yml`中的端口映射
   - 使用`kill -9 进程ID`关闭占用端口的进程

## 维护建议

1. **定期备份**
   ```bash
   # 备份数据库
   docker exec jztc-mysql sh -c 'exec mysqldump -uroot -p"root123456" go-jztc' > backup.sql
   ```

2. **日志管理**
   - 日志文件存储在`/build/logs`目录
   - 可以配置日志清理策略，避免磁盘空间耗尽
   - 建议使用`logrotate`进行日志轮转

3. **版本更新**
   - 更新前先进行备份
   - 停止现有服务: `./stop.sh`
   - 上传新版构建文件
   - 执行部署脚本: `./deploy.sh`

4. **监控建议**
   - 使用`./status.sh`查看基本状态
   - 考虑使用Prometheus+Grafana进行更全面的监控
   - 监控重点：CPU使用率、内存占用、请求响应时间 
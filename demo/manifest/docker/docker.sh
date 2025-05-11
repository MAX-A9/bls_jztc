#!/bin/bash

# 这个脚本在Docker构建前执行

# 创建必要的目录
mkdir -p ../temp/linux_amd64
mkdir -p ./temp

# 创建等待MySQL脚本
cat > ./temp/wait-for-mysql.sh << 'EOL'
#!/bin/bash

set -e

HOST="$1"
PORT="$2"
DB_USER="$3"
DB_PASS="$4"
TIMEOUT=60
RETRIES=5
RETRY_INTERVAL=3

echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [INFO] 等待数据库连接... maxRetries $RETRIES retryInterval ${RETRY_INTERVAL}s"

for i in $(seq 1 $RETRIES); do
    # 尝试连接MySQL
    if mysqladmin ping -h"$HOST" -P"$PORT" -u"$DB_USER" -p"$DB_PASS" --silent; then
        echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [INFO] 数据库连接成功"
        break
    fi
    
    # 如果是最后一次尝试，就退出
    if [ $i -eq $RETRIES ]; then
        echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [ERROR] 无法连接到数据库，超过最大重试次数"
        exit 1
    fi
    
    # 等待一段时间后重试
    echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [INFO] 无法连接到数据库，${RETRY_INTERVAL}秒后重试 ($i/$RETRIES)"
    sleep $RETRY_INTERVAL
done

# 启动应用
echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [INFO] 启动应用..."
exec "$@"
EOL

# 设置执行权限
chmod +x ./temp/wait-for-mysql.sh

# 创建启动包装脚本
cat > ./temp/start-app.sh << 'EOL'
#!/bin/bash

# 启动应用前设置环境变量
export GF_GCFG_FILE=manifest/config/config.yaml
export GF_GLOG_PATH=/app/logs

# 检查配置文件
if [ -f "/app/manifest/config/config.yaml" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [INFO] 配置文件检查通过"
else
    echo "$(date '+%Y-%m-%d %H:%M:%S.%3N') [ERROR] 配置文件不存在: /app/manifest/config/config.yaml"
    ls -la /app
    ls -la /app/manifest 2>/dev/null || echo "Manifest目录不存在"
    ls -la /app/manifest/config 2>/dev/null || echo "Config目录不存在"
    exit 1
fi

# 等待MySQL就绪
./wait-for-mysql.sh "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" ./main
EOL

# 设置执行权限
chmod +x ./temp/start-app.sh

echo "Docker构建准备完成"






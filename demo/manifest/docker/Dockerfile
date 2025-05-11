FROM loads/alpine:3.8

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
# 安装mysql客户端，用于等待mysql就绪
RUN apk add --no-cache mysql-client

# 添加资源文件
ADD resource                $WORKDIR/resource/
# 添加配置文件
ADD manifest/config         $WORKDIR/manifest/config/
# 添加SQL文件
ADD manifest/sql            $WORKDIR/manifest/sql/
# 添加主程序
ADD ./temp/linux_amd64/main $WORKDIR/main

# 创建等待MySQL脚本
RUN echo '#!/bin/sh\n\
set -e\n\
\n\
HOST="$1"\n\
PORT="$2"\n\
DB_USER="$3"\n\
DB_PASS="$4"\n\
TIMEOUT=60\n\
RETRIES=5\n\
RETRY_INTERVAL=3\n\
\n\
echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [INFO] 等待数据库连接... maxRetries $RETRIES retryInterval ${RETRY_INTERVAL}s"\n\
\n\
for i in $(seq 1 $RETRIES); do\n\
    # 尝试连接MySQL\n\
    if mysqladmin ping -h"$HOST" -P"$PORT" -u"$DB_USER" -p"$DB_PASS" --silent; then\n\
        echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [INFO] 数据库连接成功"\n\
        break\n\
    fi\n\
    \n\
    # 如果是最后一次尝试，就退出\n\
    if [ $i -eq $RETRIES ]; then\n\
        echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [ERROR] 无法连接到数据库，超过最大重试次数"\n\
        exit 1\n\
    fi\n\
    \n\
    # 等待一段时间后重试\n\
    echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [INFO] 无法连接到数据库，${RETRY_INTERVAL}秒后重试 ($i/$RETRIES)"\n\
    sleep $RETRY_INTERVAL\n\
done\n\
\n\
# 启动应用\n\
echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [INFO] 启动应用..."\n\
exec "$@"' > $WORKDIR/wait-for-mysql.sh

# 创建启动脚本
RUN echo '#!/bin/sh' > $WORKDIR/start-app.sh && \
    echo '' >> $WORKDIR/start-app.sh && \
    echo '# 启动应用前设置环境变量' >> $WORKDIR/start-app.sh && \
    echo 'export GF_GCFG_FILE=manifest/config/config.yaml' >> $WORKDIR/start-app.sh && \
    echo 'export GF_GLOG_PATH=/app/logs' >> $WORKDIR/start-app.sh && \
    echo '' >> $WORKDIR/start-app.sh && \
    echo '# 检查配置文件' >> $WORKDIR/start-app.sh && \
    echo 'if [ -f "/app/manifest/config/config.yaml" ]; then' >> $WORKDIR/start-app.sh && \
    echo '    echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [INFO] 配置文件检查通过"' >> $WORKDIR/start-app.sh && \
    echo '    echo "配置文件路径: /app/manifest/config/config.yaml"' >> $WORKDIR/start-app.sh && \
    echo '    ls -la /app/manifest/config' >> $WORKDIR/start-app.sh && \
    echo 'else' >> $WORKDIR/start-app.sh && \
    echo '    echo "$(date "+%Y-%m-%d %H:%M:%S.%3N") [ERROR] 配置文件不存在: /app/manifest/config/config.yaml"' >> $WORKDIR/start-app.sh && \
    echo '    echo "当前工作目录: $(pwd)"' >> $WORKDIR/start-app.sh && \
    echo '    ls -la /app' >> $WORKDIR/start-app.sh && \
    echo '    ls -la /app/manifest 2>/dev/null || echo "Manifest目录不存在"' >> $WORKDIR/start-app.sh && \
    echo '    ls -la /app/manifest/config 2>/dev/null || echo "Config目录不存在"' >> $WORKDIR/start-app.sh && \
    echo '    exit 1' >> $WORKDIR/start-app.sh && \
    echo 'fi' >> $WORKDIR/start-app.sh && \
    echo '' >> $WORKDIR/start-app.sh && \
    echo '# 直接启动主程序' >> $WORKDIR/start-app.sh && \
    echo 'exec ./main' >> $WORKDIR/start-app.sh

# 设置执行权限
RUN chmod +x $WORKDIR/main $WORKDIR/wait-for-mysql.sh $WORKDIR/start-app.sh

# 创建日志目录
RUN mkdir -p $WORKDIR/logs

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ["sh", "./start-app.sh"]

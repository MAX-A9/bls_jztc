#!/bin/bash

# 设置错误时停止执行
set -e

# 设置输出编码为UTF-8
export LANG=en_US.UTF-8

# 定义颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 打印带颜色的信息
print_info() {
    echo -e "${GREEN}[INFO] $1${NC}"
}

print_warn() {
    echo -e "${YELLOW}[WARN] $1${NC}"
}

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker未安装，正在安装..."
        # 安装Docker
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        rm get-docker.sh
        # 启动Docker服务
        sudo systemctl start docker
        sudo systemctl enable docker
        print_info "Docker安装完成"
    else
        print_info "Docker已安装"
    fi
}

# 检查Docker Compose是否安装
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose未安装，正在安装..."
        # 安装Docker Compose
        sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        print_info "Docker Compose安装完成"
    else
        print_info "Docker Compose已安装"
    fi
}

# 检查必要的端口是否被占用
check_ports() {
    local ports=("8000" "3306")
    for port in "${ports[@]}"; do
        if netstat -tuln | grep ":$port " > /dev/null; then
            print_error "端口 $port 已被占用，请先释放该端口"
            exit 1
        fi
    done
    print_info "端口检查通过"
}

# 创建必要的目录
create_directories() {
    local dirs=("logs" "data")
    for dir in "${dirs[@]}"; do
        if [ ! -d "../../$dir" ]; then
            mkdir -p "../../$dir"
            print_info "创建目录: ../../$dir"
        fi
    done
}

# 检查配置文件
check_config_files() {
    print_info "检查配置文件..."
    
    # 检查配置文件是否存在
    if [ ! -f "../config/config.yaml" ]; then
        print_warn "配置文件不存在，尝试从其他位置复制..."
        
        # 尝试从不同位置复制配置文件
        if [ -f "../../manifest/config/config.yaml" ]; then
            mkdir -p "../config"
            cp "../../manifest/config/config.yaml" "../config/"
            print_info "配置文件已从 ../../manifest/config/ 复制"
        elif [ -f "../../../manifest/config/config.yaml" ]; then
            mkdir -p "../config"
            cp "../../../manifest/config/config.yaml" "../config/"
            print_info "配置文件已从 ../../../manifest/config/ 复制"
        else
            print_error "无法找到配置文件，请确保配置文件存在"
            exit 1
        fi
    fi
    
    # 检查SQL脚本
    if [ ! -d "../sql" ] || [ ! "$(ls -A ../sql 2>/dev/null)" ]; then
        print_warn "SQL脚本目录不存在或为空，尝试从其他位置复制..."
        
        # 尝试从不同位置复制SQL脚本
        if [ -d "../../manifest/sql" ] && [ "$(ls -A ../../manifest/sql 2>/dev/null)" ]; then
            mkdir -p "../sql"
            cp -r ../../manifest/sql/* "../sql/"
            print_info "SQL脚本已从 ../../manifest/sql/ 复制"
        elif [ -d "../../../manifest/sql" ] && [ "$(ls -A ../../../manifest/sql 2>/dev/null)" ]; then
            mkdir -p "../sql"
            cp -r ../../../manifest/sql/* "../sql/"
            print_info "SQL脚本已从 ../../../manifest/sql/ 复制"
        else
            print_warn "无法找到SQL脚本，数据库初始化可能需要手动进行"
        fi
    fi
}

# 准备Docker环境
prepare_docker_env() {
    print_info "准备Docker环境..."
    
    # 修改Docker Compose中的相对路径为绝对路径
    CURRENT_DIR=$(pwd)
    sed -i "s|../config|${CURRENT_DIR}/../config|g" docker-compose.yml
    sed -i "s|../../logs|${CURRENT_DIR}/../../logs|g" docker-compose.yml
    sed -i "s|../sql|${CURRENT_DIR}/../sql|g" docker-compose.yml
    
    print_info "Docker环境准备完成"
}

# 停止并删除旧容器
cleanup_old_containers() {
    print_info "清理旧容器..."
    docker-compose down -v || true
}

# 构建并启动容器
build_and_start() {
    print_info "构建并启动容器..."
    docker-compose up -d --build
}

# 等待服务启动
wait_for_services() {
    print_info "等待服务启动..."
    
    # 等待MySQL服务就绪
    print_info "等待MySQL服务就绪..."
    RETRIES=30
    RETRY_COUNT=0
    
    until docker exec jztc-mysql mysqladmin ping -h localhost -u go-jztc -pBeAdwSKMbeiDpkHw --silent || [ $RETRY_COUNT -eq $RETRIES ]; do
        echo -n "."
        sleep 3
        RETRY_COUNT=$((RETRY_COUNT+1))
    done
    
    if [ $RETRY_COUNT -eq $RETRIES ]; then
        print_error "MySQL服务启动超时"
        exit 1
    fi
    
    print_info "MySQL服务已就绪"
    sleep 5
}

# 检查服务状态
check_services() {
    print_info "检查服务状态..."
    if docker ps | grep -q "jztc-app" && docker ps | grep -q "jztc-mysql"; then
        print_info "服务启动成功"
    else
        print_error "服务启动失败"
        print_info "查看应用日志："
        docker logs jztc-app
        exit 1
    fi
}

# 主函数
main() {
    print_info "开始部署..."
    
    # 检查环境
    check_docker
    check_docker_compose
    check_ports
    
    # 创建目录
    create_directories
    
    # 检查配置文件
    check_config_files
    
    # 准备Docker环境
    prepare_docker_env
    
    # 清理旧容器
    cleanup_old_containers
    
    # 构建并启动
    build_and_start
    
    # 等待服务启动
    wait_for_services
    
    # 检查服务状态
    check_services
    
    print_info "部署完成！"
    print_info "应用访问地址: http://localhost:8000"
    print_info "Swagger文档地址: http://localhost:8000/swagger"
}

# 执行主函数
main 
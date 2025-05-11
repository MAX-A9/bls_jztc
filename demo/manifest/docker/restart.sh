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

# 停止服务
stop_services() {
    print_info "正在停止服务..."
    docker-compose down
    print_info "服务已停止"
}

# 启动服务
start_services() {
    print_info "正在启动服务..."
    docker-compose up -d
    print_info "服务已启动"
}

# 等待服务启动
wait_for_services() {
    print_info "等待服务启动..."
    sleep 10
}

# 检查服务状态
check_services() {
    print_info "检查服务状态..."
    if docker ps | grep -q "jztc-app" && docker ps | grep -q "jztc-mysql"; then
        print_info "服务启动成功"
    else
        print_error "服务启动失败"
        exit 1
    fi
}

# 主函数
main() {
    print_info "开始重启服务..."
    stop_services
    start_services
    wait_for_services
    check_services
    print_info "服务重启完成"
}

# 执行主函数
main 
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
    docker-compose down -v
    print_info "服务已停止"
}

# 清理未使用的镜像和卷
cleanup() {
    print_info "清理未使用的资源..."
    docker system prune -f
    print_info "清理完成"
}

# 主函数
main() {
    print_info "开始停止服务..."
    stop_services
    cleanup
    print_info "所有服务已停止并清理完成"
}

# 执行主函数
main 
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

# 检查容器状态
check_containers() {
    print_info "容器状态："
    docker ps -a | grep -E "jztc-app|jztc-mysql"
}

# 检查服务日志
check_logs() {
    print_info "应用服务日志："
    docker logs jztc-app --tail 50
    
    print_info "数据库服务日志："
    docker logs jztc-mysql --tail 50
}

# 检查资源使用情况
check_resources() {
    print_info "资源使用情况："
    docker stats --no-stream jztc-app jztc-mysql
}

# 检查网络状态
check_network() {
    print_info "网络状态："
    docker network ls | grep jztc-network
}

# 主函数
main() {
    print_info "开始检查服务状态..."
    check_containers
    echo
    check_logs
    echo
    check_resources
    echo
    check_network
    print_info "状态检查完成"
}

# 执行主函数
main 
#!/bin/bash

function check_docker_macos() {
    # 检查 Docker 是否正在运行
    if ! docker info &> /dev/null; then
        echo "Docker 未启动，正在启动 Docker..."
        open --background -a Docker
        echo "Docker 已启动"
        sleep 10
        docker version
    else
        echo "Docker 已经在运行"
    fi
}

function check_docker_linux() {
    # 检查 Docker 是否正在运行
    if ! docker info &> /dev/null; then
        echo "Docker 未启动，正在启动 Docker..."
        systemctl start docker
        echo "Docker 已启动"
        sleep 5
        docker version
    else
        echo "Docker 已经在运行"
    fi
}

function os_type() {
    # 获取操作系统名称
    os_name=$(uname -s)
    # 判断操作系统
    if [ "$os_name" = "Darwin" ]; then
        echo "这是 macOS"
        # 在这里添加针对 macOS 的操作
        check_docker_macos
    elif [ "$os_name" = "Linux" ]; then
        echo "这是 Linux"
        # 在这里添加针对 Linux 的操作
        check_docker_linux
    else
        echo "未知操作系统"
    fi
}


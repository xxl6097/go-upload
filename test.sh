#!/bin/bash

# 检查是否提供了足够的参数
if [ $# -eq 0 ]; then
    echo "请提供路径作为参数。"
    exit 1
fi

# 使用$1获取第一个参数（路径）
path=$1

# 检查路径是否为相对路径
if [[ $path == /* ]]; then
    echo "接收到的绝对路径参数是: $path"
else
    # 转换相对路径为绝对路径
    absolute_path=$(realpath "$path")
    echo "接收到的相对路径参数是: $path"
    echo "转换为绝对路径: $absolute_path"
fi

# 在这里你可以使用$absolute_path变量进行其他操作
# 例如：ls $absolute_path 或者其他文件操作

# 结束脚本
exit 0

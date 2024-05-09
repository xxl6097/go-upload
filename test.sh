#!/bin/bash
#!/bin/bash

# 设置颜色变量
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
BLUE='\033[34m'
MAGENTA='\033[35m'
CYAN='\033[36m'
WHITE='\033[1;37m'
RESET='\033[0m'

# 使用颜色输出文本
echo -e "${RED}This is red text.${RESET}"
echo -e "${GREEN}This is green text.${RESET}"
echo -e "${YELLOW}This is yellow text.${RESET}"
echo -e "${BLUE}This is blue text.${RESET}"
echo -e "${MAGENTA}This is magenta text.${RESET}"
echo -e "${CYAN}This is cyan text.${RESET}"
echo -e "${WHITE}This is white text.${RESET}"

# 设置粗体文本
echo -e "${RED}${1m}This is red bold text.${RESET}"

# 设置下划线文本
echo -e "${BLUE}${4m}This is blue underlined text.${RESET}"
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

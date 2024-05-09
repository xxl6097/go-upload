#!/bin/bash
# fonts color
Green="\033[32m"
Red="\033[31m"
#White="\033[37m"
Yellow="\033[33m"
RedBG="\033[41;37m"
Font="\033[0m"


function log() {
#  message="$1[clink Log]: $2$3"
#  echo -e "${message}"
  local prefix="$(date '+%Y-%m-%d %H:%M:%S') [${0##*/}]"
  echo -e "$prefix${@}"
}
function info() {
  log ${Green}$1${Font}
}
function debug() {
  log "" $1 ""
}
# 日志函数
logs() {
  local prefix="$(date '+%Y-%m-%d %H:%M:%S') [${0##*/}]"
  echo "$prefix $1"
}



function test() {
    read -p "请输入来帮服务地址：" server_url
    info "来帮地址   : $server_url"
}

test
#!/bin/bash
# fonts color
Green="\033[32m"
Red="\033[31m"
Yellow="\033[33m"
RedBG="\033[41;37m"
Font="\033[0m"
# fonts color

function readpassword() {
  read -s -p "enter your password:" password
  echo "password = $password"

  case "$1" in
  -c) echo "-c-->$1,$2,$3,$4" ;;
  *) echo "-->$1,$2,$3" ;;
  esac
}

function error() {
  echo -e "${Red}$1${Font}"
}
function info() {
  echo -e "${Green}$1${Font}"
}
function warn() {
  echo -e "${Yellow}$1${Font}"
}
function errorb() {
  echo -e "${RedBG}$1${Font}"
}
function debug() {
  echo "$1"
}

function main() {
  error "hello"
  info "hello"
  warn "hello"
  errorb "hello"
  debug "hello"
}
function log() {
  message="[1Panel Log]: $1 "
  echo -e "${message}"
}
function Set_Dir() {
  if read -t 120 -p "设置 1Panel 安装目录（默认为/opt）：" TEMP_VAR; then
    if [[ "$TEMP_VAR" != "" ]]; then
      if [[ "$TEMP_VAR" != /* ]]; then
        log "输入路径不对，请重新输入"
        Set_Dir
      fi

      if [[ ! -d $TEMP_VAR ]]; then
        log "您选择的安装路径为 $TEMP_VAR"
      fi
    else
      TEMP_VAR=/opt
      log "您选择的安装路径为 $TEMP_VAR"
    fi
  else
    TEMP_VAR=/opt
    log "(设置超时，使用默认安装路径 /opt)"
  fi
  echo $TEMP_VAR
}

# shellcheck disable=SC2120
function read_arg() {
  read -p "$1" TEMP_VAR
  if [[ "$TEMP_VAR" == "" ]]; then
     read_arg $1
  fi
  echo $TEMP_VAR
}


function read_arg_timer() {
  if read -t 120 -p "$1（默认为$2)：" TEMP_VAR; then
    if [[ "$TEMP_VAR" == "" ]]; then
      TEMP_VAR=$2
    fi
  else
    TEMP_VAR=$2
  fi
  echo $TEMP_VAR
}

function test001() {
  res=$(read_arg_timer "请输入来帮地址:" "uuxia.cn")
  #  echo "----> $res"
  if [ "$res" == "" ]; then
    echo "ok"
  else
    echo "result = $res"
  fi
}

function test002() {
  res=$(read_arg "请输入来帮地址:")
  if [ "$res" == "" ]; then
    echo "ok"
  else
    echo "result = $res"
  fi
}
test002

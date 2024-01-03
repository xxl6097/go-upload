#!/bin/bash
#修改为自己的应用名称
appname=go-upload
#版本号，latest
appversion=0.0.0

function build_windows_amd64() {
    set CGO_ENABLED=0 && set GOOS=windows && set GOARCH=amd64 && go build -o ${appname}
}
function build_linux_amd64() {
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${appname}
}

function build_linux_arm64() {
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${appname}
}

function build_darwin_arm64() {
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${appname}
}

function build_images_to_hubdocker() {
    #这个地方登录一次就够了
    #docker login -u xxl6097
    docker build -t ${appname} .
    docker tag ${appname}:${appversion} xxl6097/${appname}:${appversion}
    docker buildx build --platform linux/amd64,linux/arm64 -t xxl6097/${appname}:${appversion} --push .
}

function build_images_to_conding() {
    docker login -u prdsl-1683373983040 -p ffd28ef40d69e45f4e919e6b109d5a98601e3acd clife-devops-docker.pkg.coding.net
    docker build -t ${appname} .
    docker tag ${appname}:${appversion} clife-devops-docker.pkg.coding.net/public-repository/prdsl/${appname}:${appversion}
    docker buildx build --platform linux/amd64,linux/arm64 -t clife-devops-docker.pkg.coding.net/public-repository/prdsl/${appname}:${appversion} --push .
}


function menu() {
  echo "0. 编译 Windows amd64"
  echo "1. 编译 Linux amd64"
  echo "2. 编译 Linux arm64"
  echo "3. 编译 MacOS"
  echo "4. 打包上传linux amd64/arm64镜像"
  echo "5. 打包上传linux amd64/arm64镜像-coding"
  echo "请输入编号:"
  read index

  case "$index" in
  [0]) (build_windows_amd64) ;;
  [1]) (build_linux_amd64) ;;
  [2]) (build_linux_arm64) ;;
  [3]) (build_darwin_arm64) ;;
  [4]) (build_images_to_hubdocker) ;;
  [5]) (build_images_to_conding) ;;
  *) echo "exit" ;;
  esac
}

menu

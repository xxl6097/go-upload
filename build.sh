#!/bin/bash
#修改为自己的应用名称
appname=go-upload
#版本号，latest
#appversion=0.0.0-$(date +"%Y%m%d%H%M%S")
appversion=0.0.0
isok='n'

function getversion1() {
  version=$(cat version)
  if [ "$version" = "" ]; then
    version="0.0.0"
    echo $version >version
  fi
  echo $version
}

function getversion() {
  version=$(cat version)
  if [ "$version" = "" ]; then
    version="0.0.0"
    echo $version
  else
    v3=$(echo $version | awk -F'.' '{print($3);}')
    v2=$(echo $version | awk -F'.' '{print($2);}')
    v1=$(echo $version | awk -F'.' '{print($1);}')
    if [[ $(expr $v3 \>= 9) == 1 ]]; then
      v3=0
      if [[ $(expr $v2 \>= 9) == 1 ]]; then
        v2=0
        v1=$(expr $v1 + 1)
      else
        v2=$(expr $v2 + 1)
      fi
    else
      v3=$(expr $v3 + 1)
    fi
    ver="$v1.$v2.$v3"
    echo $ver
  fi

}

function build_windows_amd64() {
#  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${appname}.exe
  docker_push_result=$(CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${appname}.exe 2>&1)
}

function build_linux_amd64() {
#  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${appname}
  docker_push_result=$(CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${appname} 2>&1)
}

function build_linux_arm64() {
#  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${appname}
  docker_push_result=$(CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${appname} 2>&1)
}

function build_darwin_arm64() {
#  CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${appname}
  docker_push_result=$(CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${appname} 2>&1)
}

function build_images_to_tencent() {
  docker login ccr.ccs.tencentyun.com --username=100016471941 -p het002402
  docker build -t ${appname} .
  docker tag ${appname}:${appversion} ccr.ccs.tencentyun.com/100016471941/${appname}:${appversion}
  docker buildx build --platform linux/amd64,linux/arm64 -t ccr.ccs.tencentyun.com/100016471941/${appname}:${appversion} --push .
}

function build_images_to_hubdocker() {
  #这个地方登录一次就够了
  docker login -u xxl6097 -p het002402
  #docker login ghcr.io --username xxl6097 --password-stdin
  docker build --build-arg ARG_VERSION="${appversion}" -t ${appname} .
  docker tag ${appname}:${appversion} xxl6097/${appname}:${appversion}
  docker buildx build --build-arg ARG_VERSION="${appversion}" --platform linux/amd64,linux/arm64 -t xxl6097/${appname}:${appversion} --push .
  #sh 'docker buildx build --platform linux/amd64,linux/arm64 -t clife-devops-docker.pkg.coding.net/public-repository/$DEPLOY_ENV/$SERVICE_NAMES:$ServiceVersion -f Dockerfile --push .'

  docker tag ${appname}:${appversion} xxl6097/${appname}:latest
  #docker buildx build --build-arg ARG_VERSION="${appversion}" --platform linux/amd64,linux/arm64 -t xxl6097/${appname}:latest --push .
  docker_push_result=$(docker buildx build --build-arg ARG_VERSION="${appversion}" --platform linux/amd64,linux/arm64 -t xxl6097/${appname}:latest --push . 2>&1)
  echo "docker pull xxl6097/${appname}:${appversion}"
  #docker run -d -p 9911:8080 --name go-raspberry xxl6097/${appname}:${appversion}
}

function build_images_to_conding() {
  docker login -u prdsl-1683373983040 -p ffd28ef40d69e45f4e919e6b109d5a98601e3acd clife-devops-docker.pkg.coding.net
  docker build -t ${appname} .
  docker tag ${appname}:${appversion} clife-devops-docker.pkg.coding.net/public-repository/prdsl/${appname}:${appversion}
#  docker buildx build --platform linux/amd64,linux/arm64 -t clife-devops-docker.pkg.coding.net/public-repository/prdsl/${appname}:${appversion} --push .
  docker_push_result=$(docker buildx build --platform linux/amd64,linux/arm64 -t clife-devops-docker.pkg.coding.net/public-repository/prdsl/${appname}:${appversion} --push . 2>&1)
  echo "docker pull clife-devops-docker.pkg.coding.net/public-repository/prdsl/${appname}:${appversion}"
}

function gomodtidy() {
  go mod tidy
}

function check_docker_macos() {
  # 检查 Docker 是否正在运行
  if ! docker info &>/dev/null; then
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
  if ! docker info &>/dev/null; then
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
  os_name=$(uname -s)
  if [ "$os_name" = "Darwin" ]; then
    check_docker_macos
  elif [ "$os_name" = "Linux" ]; then
    check_docker_linux
  else
    echo "未知操作系统"
  fi
}

function menu() {
  os_type
  echo -e "\r\n0. 编译 Windows amd64"
  echo "1. 编译 Linux amd64"
  echo "2. 编译 Linux arm64"
  echo "3. 编译 MacOS"
  echo "4. 打包多平台镜像->DockerHub"
  echo "5. 打包多平台镜像->Coding"
  echo "6. go mod tidy"
  echo "7. 打包多平台镜像->Tencent"
  echo "请输入编号:"
  read index

  appversion=$(getversion)
  echo "start===>$appversion"
  case "$index" in
  [0]) (build_windows_amd64) ;;
  [1]) (build_linux_amd64) ;;
  [2]) (build_linux_arm64) ;;
  [3]) (build_darwin_arm64) ;;
  [4]) (build_images_to_hubdocker) ;;
  [5]) (build_images_to_conding) ;;
  [6]) (gomodtidy) ;;
  [7]) (build_images_to_tencent) ;;
  *) echo "exit" ;;
  esac

  # 获取命令的退出状态码
  exit_status=$?
  # 检查退出状态码
  if [ $exit_status -eq 0 ]; then
    echo "成功"
    echo $appversion >version
  else
    echo "失败"
    echo "【$docker_push_result】"
  fi

  #  if read -t 10 -p "确定执行成功了吗:(y/n)" isok; then
  #    echo "-->$isok"
  #  else
  #    isok="y"
  #  fi
  #  if [ "$isok" = "y" ]; then
  #    echo $appversion >version
  #  fi
  rm -rf files
  git add .
  git commit -m "$appversion"
  git push

}

menu

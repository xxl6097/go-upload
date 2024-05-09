#!/bin/bash
clink_feign_arr=("https://zhyl-linying.cn/clink/api/clifemode" "https://open.api.clife.cn/clink/api/clifemode" "https://itest.clife.net/clink/api/clifemode")
clink_mqtt_conf_arr=("https://zhyl-linying.cn/v1/app" "https://api.clife.cn/v1/app" "https://pre.api.clife.cn/v1/app")
# 定义数组
function menu() {
  local options=($*)
  for ((i = 0; i < ${#options[@]}; i++)); do
    echo "$((i + 1)). ${options[$i]}"
  done
  echo "$((i + 1)). 自定义输入"
  # 获取用户输入
  read -p "请输入选项的序号: " choice

  # 验证输入是否为数字且在范围内
  if ! [[ "$choice" =~ ^[0-9]+$ ]] || ((choice < 1)) || ((choice - 1 > ${#options[@]})); then
    echo "错误的选项，请输入正确的序号。"
    menu $*
  else
    if ((choice - 1 == ${#options[@]})); then
      read -p "请输入自定义内容: " selected_option
    else
      # 根据用户选择执行相应操作
      selected_index=$((choice - 1))
      selected_option=${options[$selected_index]}
    fi
  fi
}
function save_config() {
  cat >${FRP_PATH}/${FRP_NAME}.toml <<EOF
[lonbon]
# 来邦地址，填之前，请先ping一下地址是否通
LONBON_SERVER_URL= uuxia.cn
# 来邦端口
LONBON_SERVER_POST = 6302
# 来邦机构ID
LONBON_ORGID = 245826234907136539
# 来邦账户，来帮私有管理地址示例：http://192.168.x.x:8080
LONBON_API_USERNAME = clink
# 来邦密码
LONBON_API_PASSWORD = Het@1234
# 来邦api地址，同样，填之前判断一下能否正常访问
LONBON_FEIGN_URL = http://uuxia.cn:6303

[clife]
# 地址填写之前，测试一下能否正常访问
# 临颍地址：https://zhyl-linying.cn/clink/api/clifemode
# 数联生产环境：https://open.api.clife.cn/clink/api/clifemode
# 数联itest环境：https://itest.clife.net/clink/api/clifemode
CLINK_FEIGN_URL=https://zhyl-linying.cn/clink/api/clifemode
# 地址填写之前，测试一下能否正常访问
# 临颍地址：https://zhyl-linying.cn/v1/app
# 数联生产环境：https://api.clife.cn/v1/app
# 数联itest环境：https://pre.api.clife.cn/v1/app
CLIFE_MQTT_CONFIG_URL=https://zhyl-linying.cn/v1/app
EOF
}
function testmenu() {
  arr=("https://zhyl-linying.cn/clink/api/clifemode" "https://open.api.clife.cn/clink/api/clifemode" "https://itest.clife.net/clink/api/clifemode")
  echo "请输入："
  menu ${arr[*]}
  echo "->$selected_option"
}
function main() {
    echo "请选择clink_feign_url地址："
    menu ${clink_feign_arr[*]}
    clink_feign_url=$selected_option
    echo -e "\r\n请选择clink_mqtt_url地址："
    menu ${clink_mqtt_conf_arr[*]}
    clink_mqtt_url=$selected_option
    echo "-->$clink_feign_url    $clink_mqtt_url"
}
main
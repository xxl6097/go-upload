#!/bin/bash

read -s -p "enter your password:" password
echo "password = $password"

case "$1" in
    -c) echo "-c-->$1,$2,$3,$4";;
     *) echo "-->$1,$2,$3" ;;
esac


#!/bin/bash
Var1=200

# 定义函数
function fun1(){
        Var1=121
}

#执行函数
fun1
echo after func Var1 is $0:$Var1
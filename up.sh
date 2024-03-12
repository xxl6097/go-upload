#!/bin/bash
cmd="curl "
for arg in "$@"; do
  if [[ $arg == /* ]]; then
      cmd+="-F \"file=@$arg\" "
  else
      absolute_path=$(realpath "$arg")
      cmd+="-F \"file=@$absolute_path\" "
  fi
done
cmd+="-F \"token=het002402\" http://uuxia.cn:8087/upload"
echo "运行命令：$cmd"
eval $cmd


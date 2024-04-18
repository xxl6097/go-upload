#!/bin/bash
cmd="curl"
header="-H \"Authorization: het002402\""
host="http://localhost:4444/upload"
files=""
for arg in "$@"; do
  if [[ $arg == /* ]]; then
      files+="-F \"file=@$arg\" "
  else
      absolute_path=$(realpath "$arg")
      files+="-F \"file=@$absolute_path\" "
  fi
done
cmd="curl $header $files$host"
echo "$cmd"
eval $cmd
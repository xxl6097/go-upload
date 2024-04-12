#!/bin/bash
cmd="curl "
for arg in "$@"; do
  case $arg in
  token*) cmd+="-F \"$arg\" " ;;
  *) if [[ $arg == /* ]]; then
    if [ -n "$arg" ]; then
      cmd+="-F \"file=@$arg\" "
    fi
  else
    absolute_path=$(realpath "$arg")
    if [ -n "$absolute_path" ]; then
      cmd+="-F \"file=@$absolute_path\" "
    fi
  fi ;;
  esac
done
if [[ $cmd != *"token"* ]]; then
  read -s -p "enter token:" token
  echo "token = $token"
  cmd+="-F \"token=$token\""
fi
cmd+=" http://localhost:4444/upload"
echo "$cmd"
eval $cmd

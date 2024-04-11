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
echo -e "\033[31mplease input token:\033[0m"
read token
cmd+="-F \"token=$token\" http://localhost:4444/upload"
echo "run cmd: $cmd"
eval $cmd

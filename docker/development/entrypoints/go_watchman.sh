#!/usr/bin/env bash

cd "$BACKEND_DIR"

server_pid="$(pgrep server)"

if [[ -n "$server_pid" ]]; then
  kill -KILL "$server_pid"
  while kill -0 "$server_pid" 2>/dev/null; do sleep 1; done
fi

exec 1>>/proc/1/fd/1
exec 2>>/proc/1/fd/2

echo -e "\nBuilding Server"
gb build server

if [[ $? -ne 0 ]]; then
  echo -e "\nBuild Failed not running"
  exit 1
fi

echo -e "\nRunning Server"
"$BACKEND_BINARY" &

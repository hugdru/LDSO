#!/usr/bin/env bash

dockerize -wait tcp://mongodb:27017 -wait tcp://redis:6379 -timeout 30s

cd "$BACKEND_DIR"

echo 'Fetching libraries'
gb vendor restore
watchman -n -j <<EOF
["trigger", "$BACKEND_DIR/src", {
  "name": "go_assets",
  "command": ["watchman_command.sh"],
  "append_files": false
}]
EOF
tail -f /dev/null

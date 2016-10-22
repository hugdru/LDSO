#!/usr/bin/env bash

gb vendor restore
gb build server
"$BACKEND_BINARY" &
watchman -f -j <<EOF
["trigger", "src", {
  "name": "go_assets",
  "command": ["watchman_command.sh"],
  "append_files": false
}]
EOF

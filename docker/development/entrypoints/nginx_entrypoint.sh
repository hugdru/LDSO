#!/usr/bin/env bash

cd "$FRONTEND_DIR"
npm install && watchman -j <<"EOF"
["trigger", ".", {
  "name": "angular_assets",
  "command": ["ng", "build"],
  "append_files": false
}]
EOF

#!/usr/bin/env bash

server_pid="$(pgrep server)"

kill -KILL "$server_pid"
while kill -0 "$server_pid" 2>/dev/null; do sleep 1; done
gb build server
"$BACKEND_BINARY" &

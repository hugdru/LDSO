#!/usr/bin/env bash

nginx

cd "$FRONTEND_DIR"
# ng build -dev --watch true
ng serve --port 4200 --host 0.0.0.0

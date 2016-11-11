#!/usr/bin/env bash

nginx

cd "$FRONTEND_DIR"
# ng build --env="$BUILD" --watch true
npm install
ng serve --env="$BUILD" --port 4200 --host 0.0.0.0

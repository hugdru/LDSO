#!/usr/bin/env bash

dockerize -wait tcp://mongodb:27017 -wait tcp://redis:6379 -timeout 10s ./bin/server

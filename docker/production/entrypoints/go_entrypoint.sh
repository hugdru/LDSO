#!/usr/bin/env bash

dockerize -wait tcp://postgres:5432 -wait tcp://redis:6379 -timeout 10s ./bin/server

#!/usr/bin/env bash

dockerize -wait tcp://postgres:5432 -timeout 10s ./bin/server

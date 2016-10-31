#!/usr/bin/env bash

curl -X POST -d "{\"name\":\"Estacionamento\", \"weight\":30, \"sub_groups\":[]}" "http://localhost:8080/setGroup"

curl -X GET -d "http://localhost:8080/getAllGroups"
#!/usr/bin/env bash

cd "${0%/*}"

curl -X GET "http://localhost:8080/getAllGroups"

curl -X POST -d "{\"name\":\"Acessos\"}" "http://localhost:8080/getGroup"

curl -X POST -d "{\"name\":\"Estacionamento\", \"weight\":30, \"sub_groups\":[]}" "http://localhost:8080/setGroup"


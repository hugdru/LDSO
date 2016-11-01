#!/usr/bin/env bash

cd "${0%/*}"

#curl -X GET "http://localhost:8080/getAllGroups"

#curl -X GET "http://localhost:8080/getGroup?name=Percurso%20Interior"

#curl -X POST -d "{\"name\":\"Estacionamento\", \"weight\":30, \"sub_groups\":[]}" "http://localhost:8080/setGroup"

curl -X POST -d "{\"name\":\"Paralelo\", \"weight\":10, \"criteria\":[]}" "http://localhost:8080/setSubGroup?name=Percurso%20Interior"

#curl -X GET "http://localhost:8080/getGroup?name=Percurso%20Interior"

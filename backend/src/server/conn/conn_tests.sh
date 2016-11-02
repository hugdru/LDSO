#!/usr/bin/env bash

cd "${0%/*}"

curl -X GET "http://localhost:8080/getAllGroups"

curl -X POST -d "{\"_id\":5, \"name\":\"Coisas\", \"weight\":30}" "http://localhost:8080/setMainGroup"

curl -X GET "http://localhost:8080/getMainGroup?name=Coisas"

curl -X POST -d "{\"_id\":3, \"name\":\"Paralelo\", \"weight\":10, \"main_group\":2}" "http://localhost:8080/setSubGroup"

curl -X GET "http://localhost:8080/getSubGroup?name=Paralelo"

curl -X POST -d "{\"_id\":3, \"name\":\"xpto\", \"weight\":5, \"legislation\":\"\", \"sub_group\":2}" "http://localhost:8080/setCriterion"

curl -X GET "http://localhost:8080/getCriterion?name=xpto"

curl -X POST -d "{\"_id\":9, \"name\":\"fisica\", \"weight\":5, \"criterion\":3}" "http://localhost:8080/setAccessibility"

curl -X GET "http://localhost:8080/getAccessibility?name=fisica"

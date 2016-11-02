#!/usr/bin/env bash

cd "${0%/*}"

curl -X GET "http://localhost:8080/getAllMainGroups"

curl -X POST -d "{\"_id\":5, \"name\":\"Coisas\", \"weight\":30}" "http://localhost:8080/setMainGroup"

curl -X GET "http://localhost:8080/getMainGroup?tag=name&type=string&value=Coisas"

curl -X GET "http://localhost:8080/updateMainGroup?_id=5&tag=weight&type=int&value=20"

curl -X GET "http://localhost:8080/updateMainGroup?_id=5&tag=name&type=string&value=Cenas"

curl -X GET "http://localhost:8080/getOneMainGroup?tag=_id&type=int&value=5"

curl -X GET "http://localhost:8080/getAllSubGroups"

curl -X POST -d "{\"_id\":3, \"name\":\"Paralelo\", \"weight\":10, \"main_group\":2}" "http://localhost:8080/setSubGroup"

curl -X GET "http://localhost:8080/getOneSubGroup?tag=name&type=string&value=Paralelo"

curl -X GET "http://localhost:8080/updateSubGroup?_id=3&tag=name&type=string&value=Paralelinho"

curl -X GET "http://localhost:8080/getOneSubGroup?tag=name&type=string&value=Paralelinho"

curl -X GET "http://localhost:8080/getAllCriteria"

curl -X POST -d "{\"_id\":3, \"name\":\"xpto\", \"weight\":5, \"legislation\":\"\", \"sub_group\":2}" "http://localhost:8080/setCriterion"

curl -X GET "http://localhost:8080/getOneCriterion?tag=name&type=string&value=xpto"

curl -X GET "http://localhost:8080/updateCriterion?_id=3&tag=name&type=string&value=yqup"

curl -X GET "http://localhost:8080/getOneCriterion?tag=name&type=string&value=yqup"

curl -X GET "http://localhost:8080/getAllAccessibilities"

curl -X POST -d "{\"_id\":9, \"name\":\"fisica\", \"weight\":5, \"criterion\":3}" "http://localhost:8080/setAccessibility"

curl -X GET "http://localhost:8080/getAccessibilities?tag=name&type=string&value=fisica"

curl -X GET "http://localhost:8080/updateAccessibility?_id=9&tag=weight&type=int&value=5"

curl -X GET "http://localhost:8080/getOneAccessibility?tag=_id&type=int&value=9"

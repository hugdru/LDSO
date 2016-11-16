#!/usr/bin/env bash
cd "${0%/*}"

./../../../../database/populate.sh

EXE="curl -X"
URL="http://go1:8080/"

PostTests=( \
    "{\"name\":\"Coisas\",\"weight\":30}" \
    "mainGroups" \
#
    "{\"name\":\"Paralelo\",\"weight\":10,\"main_group\":2}" \
    "subGroups" \
#
    "{\"name\":\"xpto\",\"weight\":5,\"legislation\":\"\",\"sub_group\":2}" \
    "criteria" \
#
    "{\"name\":\"fisica\",\"weight\":15,\"criterion\":3}" \
    "accessibilities" \
#
    '{"property": 1,"rating": 30, "criteria": [{"criterion": 1, "value": 30},{"criterion": 2, "value": 20}]}' \
    "audits" \
)

PutTests=( \
    "{\"_id\":5,\"name\":\"Coisas\",\"weight\":20}" \
    "mainGroups?_id=5" \
#
    "{\"_id\":5,\"name\":\"Cenas\",\"weight\":20}" \
    "mainGroups?_id=5" \
#
    "{\"_id\":3,\"name\":\"Paralelinho\",\"weight\":10,\"main_group\":2}" \
    "subGroups?_id=3" \
#
    "{\"_id\":3,\"name\":\"yqup\",\"weight\":5,\"legislation\":\"\",\"sub_group\":2}" \
    "criteria?_id=3" \
#
    "{\"_id\":9,\"name\":\"fisica\",\"weight\":5,\"criterion\":3}" \
    "accessibilities?_id=9" \
)

GetTests1=( \
    "mainGroups/find?tag=_id&type=int&value=5" \
    "subGroups/find?tag=name&type=string&value=Paralelinho" \
    "criteria/find?tag=name&type=string&value=yqup" \
    "accessibilities?tag=name&type=string&value=fisica" \
    "mainGroups" \
    "subGroups" \
    "criteria" \
    "accessibilities" \
    "properties" \
    "audits" \
)

DeleteTests=(\
    "mainGroups?_id=1" \
    "subGroups" \
)

GetTests2=( \
    "mainGroups" \
    "subGroups" \
    "criteria" \
    "accessibilities" \
)


for (( t=0; t<${#PostTests[@]}; t=t+2))
do
  $EXE POST -d "${PostTests[$t]}" $URL${PostTests[$t+1]}
done

for (( t=0; t<${#PutTests[@]}; t=t+2))
do
  $EXE PUT -d "${PutTests[$t]}" $URL${PutTests[$t+1]}
done

for t in ${GetTests1[@]}
do
  $EXE GET $URL$t
done

#for t in ${DeleteTests[@]}
#do
# $EXE DELETE $URL$t
#done
#
#for t in ${GetTests2[@]}
#do
# $EXE GET $URL$t
#done

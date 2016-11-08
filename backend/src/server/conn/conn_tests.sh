#!/usr/bin/env bash
cd "${0%/*}"

./../../../../database/populate.sh

EXE="curl -X"
URL="http://go1:8080/"

PostTests=( \
		"{\"_id\":5,\"name\":\"Coisas\",\"weight\":30}" \
		"mainGroups" \
#
		"{\"_id\":3,\"name\":\"Paralelo\",\"weight\":10,\"main_group\":2}" \
		"subGroups" \
#
		"{\"_id\":3,\"name\":\"xpto\",\"weight\":5,\"legislation\":\"\",\"sub_group\":2}" \
		"criteria" \
#
		"{\"_id\":9,\"name\":\"fisica\",\"weight\":5,\"criterion\":3}" \
		"accessibilities" \
)

PutTests=(\
		"mainGroups?_id=5&tag=weight&type=int&value=20" \
		"mainGroups?_id=5&tag=name&type=string&value=Cenas" \
		"subGroups?_id=3&tag=name&type=string&value=Paralelinho" \
		"accessibilities?_id=9&tag=weight&type=int&value=5" \
		"criteria?_id=3&tag=name&type=string&value=yqup" \
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
	$EXE POST -d ${PostTests[$t]} $URL${PostTests[$t+1]}
done

for t in ${PutTests[@]}
do
	$EXE PUT $URL$t
done

for t in ${GetTests1[@]}
do
	$EXE GET $URL$t
done

for t in ${DeleteTests[@]}
do
	$EXE DELETE $URL$t
done

for t in ${GetTests2[@]}
do
	$EXE GET $URL$t
done

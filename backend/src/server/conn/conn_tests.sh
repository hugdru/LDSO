#!/usr/bin/env bash
cd "${0%/*}"

./../../../../database/populate.sh

EXE="curl -X"
URL="http://go1:8080/"

PostTests=( \
		"{\"_id\":5,\"name\":\"Coisas\",\"weight\":30}" \
		"mainGroup" \
#
		"{\"_id\":3,\"name\":\"Paralelo\",\"weight\":10,\"main_group\":2}" \
		"subGroup" \
#
		"{\"_id\":3,\"name\":\"xpto\",\"weight\":5,\"legislation\":\"\",\"sub_group\":2}" \
		"criterion" \
#
		"{\"_id\":9,\"name\":\"fisica\",\"weight\":5,\"criterion\":3}" \
		"accessibility" \
)

PutTests=(\
		"mainGroup?_id=5&tag=weight&type=int&value=20" \
		"mainGroup?_id=5&tag=name&type=string&value=Cenas" \
		"subGroup?_id=3&tag=name&type=string&value=Paralelinho" \
		"accessibility?_id=9&tag=weight&type=int&value=5" \
		"criterion?_id=3&tag=name&type=string&value=yqup" \
)

GetTests1=( \
		"mainGroup/find?tag=name&type=string&value=Cenas" \
		"mainGroup/find?tag=_id&type=int&value=5" \
		"subGroup/find?tag=name&type=string&value=Paralelinho" \
		"criterion/find?tag=name&type=string&value=yqup" \
		"accessibility?tag=name&type=string&value=fisica" \
		"mainGroup" \
		"subGroup" \
		"criterion" \
		"accessibility" \
)

DeleteTests=(\
		"mainGroup?_id=1" \
)

GetTests2=( \
		"mainGroup" \
		"subGroup" \
		"criterion" \
		"accessibility" \
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

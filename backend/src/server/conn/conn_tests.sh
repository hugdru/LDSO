#!/usr/bin/env bash
EXE="curl -X"
URL="http://go1:8080/"
PostTests=( \
		"{\"_id\":5,\"name\":\"Coisas\",\"weight\":30}" \
		"setMainGroup" \
#
		"{\"_id\":3,\"name\":\"Paralelo\",\"weight\":10,\"main_group\":2}" \
		"setSubGroup" \
#
		"{\"_id\":3,\"name\":\"xpto\",\"weight\":5,\"legislation\":\"\",\"sub_group\":2}" \
		"setCriterion" \
#
		"{\"_id\":9,\"name\":\"fisica\",\"weight\":5,\"criterion\":3}" \
		"setAccessibility" \
)
GetTests=(\
		"getOneMainGroup?tag=name&type=string&value=Coisas" \
		"updateMainGroup?_id=5&tag=weight&type=int&value=20" \
		"updateMainGroup?_id=5&tag=name&type=string&value=Cenas" \
		"getOneMainGroup?tag=_id&type=int&value=5" \
		"getOneSubGroup?tag=name&type=string&value=Paralelo" \
		"updateSubGroup?_id=3&tag=name&type=string&value=Paralelinho" \
		"getOneSubGroup?tag=name&type=string&value=Paralelinho" \
		"getOneCriterion?tag=name&type=string&value=xpto" \
		"updateCriterion?_id=3&tag=name&type=string&value=yqup" \
		"getOneCriterion?tag=name&type=string&value=yqup" \
		"getAccessibilities?tag=name&type=string&value=fisica" \
		"updateAccessibility?_id=9&tag=weight&type=int&value=5" \
		"getAllMainGroups" \
		"getAllSubGroups" \
		"getAllCriteria" \
		"getAllAccessibilities" \
		"removeMainGroup?_id=1" \
		"getAllMainGroups" \
		"getAllSubGroups" \
		"getAllCriteria" \
		"getAllAccessibilities" \
	)
for (( t=0; t<${#PostTests[@]}; t=t+2))
do
	$EXE POST -d ${PostTests[$t]} $URL${PostTests[$t+1]}
done
for t in ${GetTests[@]}
do
	$EXE GET $URL$t
done

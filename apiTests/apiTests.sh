#!/usr/bin/env bash

set -e

cd "${0%/*}"

database_script="../database/pgr"
specs_dir="./specs"

EXE="curl -X"
URL="http://api.lp4adev.tk:8080"

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NOCOLOR='\033[0m'

typeset -A gets
gets=(
  ['countries']='countries.spec'
  ['countries/1']='country1.spec'
  ['properties/1']='property1.spec'
  ['addresses/1']='address1.spec'
  ['clients/1']='client1.spec'
  ['templates']='templates.spec'
  ['templates?id=2']='templatesId2.spec'
  ['templates?id=2&name=O+meu+segundo+template']='templatesId2Name.spec'
  ['templates?name=O+meu+primeiro+template']='templatesName1.spec'
  ['templates/1']='template1.spec'
  ['templates/2']='template2.spec'
  ['maingroups']='maingroups.spec'
  ['maingroups?idTemplate=1']='maingroupsIdTemplate1.spec'
  ['maingroups?idTemplate=1&name=T1-M2']='maingroupsIdTemplate1Name.spec'
  ['maingroups/1']='maingroup1.spec'
  ['maingroups/2']='maingroup2.spec'
  ['subgroups']='subgroups.spec'
  ['subgroups?idMaingroup=2']='subgroupsIdMaingroup2.spec'
  ['subgroups/1']='subgroup1.spec'
  ['subgroups/2']='subgroup2.spec'
  ['legislations']='legislations.spec'
  ['legislations?name=O+banco']='legislationsName1.spec'
  ['legislations/1']='legislation1.spec'
  ['legislations/2']='legislation2.spec'
  ['criteria']='criteria.spec'
  ['criteria?idSubgroup=5']='criteriaIdSubgroup5.spec'
  ['criteria/1']='criterion1.spec'
  ['criteria/2']='criterion2.spec'
  ['auditors/1']='auditors1.spec'

)

main() {

  if [[ "$1" != noinit ]]; then
    "$database_script" init
  fi
  "$database_script" ddl
  "$database_script" examples
  bail=false

  get || bail=true


  if [[ "$bail" == true ]]; then
    echo -e "${RED}TESTS FAILED!"
    exit 1
  fi

  exit 0
}

get() {
  echo -e "${BLUE}Get tests"

  local success=true
  local n_tests=${#gets[@]}

  local counter=1
  for resource in "${!gets[@]}"; do
    local spec_file="$specs_dir/${gets[$resource]}"
    if [[ ! -f "$spec_file" ]]; then
      echo -e "\t${RED}$resource \t\t $counter/$n_tests \t NO SPEC!"
      success=false
    else
      local output="$($EXE GET "$URL/$resource" 2>/dev/null)"
      local spec=$(<"$spec_file")
      if [[ "$output" == "$spec" ]]; then
        echo -e "\t${GREEN}$resource \t\t $counter/$n_tests"
      else
        echo -e "\t${RED}$resource \t\t $counter/$n_tests"
        success=false
      fi
    fi
    ((counter++))
  done

  if [[ "$success" == true ]]; then
    return 0
  fi
  return -1
}

main "$@"

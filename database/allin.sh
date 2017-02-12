#!/usr/bin/env bash

###
# This script concatenates all the files in sql_dir_path
# recursively and by lexical order.
# Then it prints the result to stdout.
# The depth can be controlled by passing an argument,
# with no argument there is no limit to the depth.
#
# Examples:
# ./allin.sh
# ./allin.sh 1
# ./allin.sh n
###

sql_dir_path="./sql"
script_dir="${0%/*}"

cd "$script_dir"

if [[ $# -gt 1 ]]; then
  echo "$0 [number]"
  exit 1
fi

if [[ $# -eq 1 ]]; then
  re='^[0-9]+$'
  if ! [[ "$1" =~ $re ]]; then
    1>&2 echo "error: must be a number"
    exit 1
  fi
  if [[ $1 -eq 0 ]]; then
    1>&2 echo "error: must not be zero"
    exit 1
  fi
  depth="-maxdepth $1"
fi

sql_files=()
echo "-- Concatenated files"
while IFS= read -r -d $'\0' sql_file; do
    echo "-- $sql_file"
    sql_files+=("$sql_file")
done < <(find "$sql_dir_path" $depth -iname "*.sql" -print0 | sort -z -V)

echo

cat "${sql_files[@]}"

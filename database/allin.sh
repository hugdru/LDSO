#!/usr/bin/env bash

cd "${0%/*}"

if [[ $# -ne 1 ]]; then
  echo "$0 file"
  exit 1
fi

sqls=()
while IFS= read -r -d $'\0' sql; do
    sqls+=("$sql")
done < <(find ./sql -name "*.sql" -print0 | sort -z -V)

cat "${sqls[@]}" > $1

#!/usr/bin/env bash

cd "${0%/*}"

DB=Places4All
collections=(groups property)

mongo $DB --eval "db.dropDatabase()"
for coll in ${collections[@]}
do
	mongoimport --db $DB --collection $coll --file $coll.json
done

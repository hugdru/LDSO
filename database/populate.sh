#!/usr/bin/env bash

cd "${0%/*}"

DB=Places4All
collections=(main_group sub_group criterion accessibility property)

mongo $DB --eval "db.dropDatabase()"

mongo $DB --eval "db.main_group.ensureIndex({name: 1}, {unique: true})"
#mongo $DB --eval "db.main_group.sub_groups.ensureIndex({name: 1}, {unique: true})" # this does not work!!!
mongo $DB --eval "db.property.ensureIndex({name: 1}, {unique: true})"

for coll in ${collections[@]}
do
	mongoimport --db $DB --collection $coll --file $coll.json
done

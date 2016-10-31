DB=Places4All
collections=(property groups)

mongo $DB --eval "db.dropDatabase()"
for coll in ${collections[@]}
do
	mongoimport --db $DB --collection $coll --file $coll.json
done

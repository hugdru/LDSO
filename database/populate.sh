DB=Places4All
collections=(property groups)

for coll in ${collections[@]}
do
	mongoimport --db $DB --collection $coll --file $coll.json
done

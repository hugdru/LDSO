package main

import (
	"net/http"
	"server/db"
	"server/conn"
	"server/data"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

func main() {
	var coll map[string]*mgo.Collection
	coll = make(map[string]*mgo.Collection)

	// TODO: Find a better way to make the index unique on first run.
	// Should be done in a json not on go.

	first_run := false

	session := db.StartConn("localhost:27017")
	defer db.CloseConn(session)

	if !db.ExistsCollections(session, "Places4All") {
		first_run = true
	}

	coll["groups"] = db.GetCollection(session, "Places4All", "group")
	coll["groups_set"] = db.GetCollection(session, "Places4All", "group_set")
	coll["property"] = db.GetCollection(session, "Places4All", "property")

	if first_run == true {
		db.EnsureUnique(coll["groups"], "name")
		db.EnsureUnique(coll["property"], "name")
	}

	populate(coll);

	http.HandleFunc("/groups", conn.GetHandlerGroup(coll["groups"]))
	http.HandleFunc("/property", conn.GetHandlerProperty(coll["property"]))

	http.ListenAndServe(":8080", nil)
}

/*
 * Test data start
 * Delete latter
 */
func populate(coll map[string]*mgo.Collection) {
	prop1 := data.Property{"Casa", data.Owner{"Joao"}, "pic1.png"}
	prop2 := data.Property{"Hotel Sunny", data.Owner{"Carlos"}, "pic2.png"}

	db.Insert(coll["property"], &prop1, &prop2)

	group1 := data.Group{"Acessos", 40}
	group2 := data.Group{"Percurso Exterior", 20}
	group3 := data.Group{"Percurso Interior", 20}
	group4 := data.Group{"Bens e Serviços", 20}

	db.Insert(coll["groups"], &group1, &group2, &group3, &group4)

	access1, access2 := data.Accessibility("Type a"),
			data.Accessibility("Type b")

	crit1 := data.Criterion{"Rampa", access1, true}
	crit2 := data.Criterion{"Porta", access1, false}
	crit3 := data.Criterion{"Altura", access2, true}
	crit4 := data.Criterion{"Entrada", access2, false}
	crit5 := data.Criterion{"WC", access1, true}
	crit6 := data.Criterion{"Cadeira", access2, true}

	sub1 := data.Sub_Group{"A", 30, nil}
	sub2 := data.Sub_Group{"B", 70, nil}

	sub1.SetCriteria(crit1, crit2, crit3)
	sub2.SetCriteria(crit4, crit5, crit6)

	set1 := data.Group_Set{group1, nil}
	set1.SetSubs(sub1)

	set2 := data.Group_Set{group2, nil}
	set2.SetSubs(sub2)

	db.Insert(coll["groups_set"], &set1, &set2)
}
/*
 * Test data end
 */

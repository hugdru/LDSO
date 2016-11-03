package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"server/conn"
	"server/db"
	// "gopkg.in/mgo.v2/bson"
)

func main() {
	var coll map[string]*mgo.Collection
	coll = make(map[string]*mgo.Collection)

	// TODO: Find a better way to make the index unique on first run.
	// Should be done in a json not on go.

	first_run := false

	session := db.StartConn("mongodb:27017")
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

	http.HandleFunc("/groups", conn.GetHandlerGroup(coll["groups"]))
	http.HandleFunc("/property", conn.GetHandlerProperty(coll["property"]))

	http.ListenAndServe(":8080", nil)
}

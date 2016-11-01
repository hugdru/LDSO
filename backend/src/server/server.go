package main

import (
	"net/http"
	"server/db"
	"server/conn"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

func main() {
	var coll map[string]*mgo.Collection
	coll = make(map[string]*mgo.Collection)

	// TODO: Find a better way to make the index unique on first run.
	// Should be done in a json not on go.

//	first_run := false

	session := db.StartConn("localhost:27017")
	defer db.CloseConn(session)

//	if !db.ExistsCollections(session, "Places4All") {
//		first_run = true
//	}

	coll["main_group"] = db.GetCollection(session, "Places4All", "main_group")
	coll["property"] = db.GetCollection(session, "Places4All", "property")

//	if first_run == true {
//		db.EnsureUnique(coll["main_group"], "name")
//		db.EnsureUnique(coll["property"], "name")
//	}

	http.HandleFunc("/getAllGroups", conn.GetAllGroups(coll["main_group"]))
	http.HandleFunc("/getGroup", conn.GetGroup(coll["main_group"]))
	http.HandleFunc("/setGroup", conn.SetGroup(coll["main_group"]))
	http.HandleFunc("/setSubGroup", conn.SetSubGroup(coll["main_group"]))
	http.HandleFunc("/property", conn.GetHandlerProperty(coll["property"]))

	http.ListenAndServe(":8080", nil)
}

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

	// TODO: Find a better way to make the index unique on first run. DONE?
	// Should be done in a json not on go.

//	first_run := false

	session := db.StartConn("localhost:27017")
	defer db.CloseConn(session)

//	if !db.ExistsCollections(session, "Places4All") {
//		first_run = true
//	}

	coll["main_group"] = db.GetCollection(session, "Places4All", "main_group")
	coll["sub_group"] = db.GetCollection(session, "Places4All", "sub_group")
	coll["criterion"] = db.GetCollection(session, "Places4All", "criterion")
	coll["accessibility"] = db.GetCollection(session, "Places4All", "accessibility")
	coll["property"] = db.GetCollection(session, "Places4All", "property")

//	if first_run == true {
//		db.EnsureUnique(coll["main_group"], "name")
//		db.EnsureUnique(coll["property"], "name")
//	}

	http.HandleFunc("/getAllMainGroups", conn.GetAll(coll["main_group"]))
	http.HandleFunc("/getMainGroups", conn.Get(coll["main_group"]))
	http.HandleFunc("/getOneMainGroup", conn.GetOne(coll["main_group"]))
	http.HandleFunc("/setMainGroup", conn.Set(coll["main_group"]))
	http.HandleFunc("/updateMainGroup", conn.Update(coll["main_group"]))
	http.HandleFunc("/getAllSubGroups", conn.GetAll(coll["sub_group"]))
	http.HandleFunc("/getSubGroups", conn.Get(coll["sub_group"]))
	http.HandleFunc("/getOneSubGroup", conn.GetOne(coll["sub_group"]))
	http.HandleFunc("/setSubGroup", conn.Set(coll["sub_group"]))
	http.HandleFunc("/updateSubGroup", conn.Update(coll["sub_group"]))
	http.HandleFunc("/getAllCriteria", conn.GetAll(coll["criterion"]))
	http.HandleFunc("/getCriteria", conn.Get(coll["criterion"]))
	http.HandleFunc("/getOneCriterion", conn.GetOne(coll["criterion"]))
	http.HandleFunc("/setCriterion", conn.Set(coll["criterion"]))
	http.HandleFunc("/updateCriterion", conn.Update(coll["criterion"]))
	http.HandleFunc("/getAllAccessibilities", conn.GetAll(coll["accessibility"]))
	http.HandleFunc("/getAccessibilities", conn.Get(coll["accessibility"]))
	http.HandleFunc("/getOneAccessibility", conn.GetOne(coll["accessibility"]))
	http.HandleFunc("/setAccessibility", conn.Set(coll["accessibility"]))
	http.HandleFunc("/updateAccessibility", conn.Update(coll["accessibility"]))
	http.HandleFunc("/property", conn.GetHandlerProperty(coll["property"]))

	http.ListenAndServe(":8080", nil)
}

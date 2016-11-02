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
	coll["sub_group"] = db.GetCollection(session, "Places4All", "sub_group")
	coll["criterion"] = db.GetCollection(session, "Places4All", "criterion")
	coll["accessibility"] = db.GetCollection(session, "Places4All", "accessibility")
	coll["property"] = db.GetCollection(session, "Places4All", "property")

//	if first_run == true {
//		db.EnsureUnique(coll["main_group"], "name")
//		db.EnsureUnique(coll["property"], "name")
//	}

	http.HandleFunc("/getAllGroups", conn.GetAllMainGroups(coll["main_group"]))
	http.HandleFunc("/getMainGroup", conn.GetMainGroup(coll["main_group"]))
	http.HandleFunc("/setMainGroup", conn.SetMainGroup(coll["main_group"]))
	http.HandleFunc("/updateMainGroup", conn.UpdateMainGroup(coll["main_group"]))
	http.HandleFunc("/getAllSubGroups", conn.GetAllSubGroups(coll["sub_group"]))
	http.HandleFunc("/getSubGroup", conn.GetSubGroup(coll["sub_group"]))
	http.HandleFunc("/setSubGroup", conn.SetSubGroup(coll["sub_group"]))
	http.HandleFunc("/updateSubGroup", conn.UpdateSubGroup(coll["sub_group"]))
	http.HandleFunc("/getAllCriteria", conn.GetAllCriteria(coll["criterion"]))
	http.HandleFunc("/getCriterion", conn.GetCriterion(coll["criterion"]))
	http.HandleFunc("/setCriterion", conn.SetCriterion(coll["criterion"]))
	http.HandleFunc("/updateCriterion", conn.UpdateCriterion(coll["criterion"]))
	http.HandleFunc("/getAllAccessibilities", conn.GetAllAccessibilities(coll["accessibility"]))
	http.HandleFunc("/getAccessibility", conn.GetAccessibility(coll["accessibility"]))
	http.HandleFunc("/setAccessibility", conn.SetAccessibility(coll["accessibility"]))
	http.HandleFunc("/updateAccessibility", conn.UpdateAccessibility(coll["accessibility"]))
	http.HandleFunc("/property", conn.GetHandlerProperty(coll["property"]))

	http.ListenAndServe(":8080", nil)
}

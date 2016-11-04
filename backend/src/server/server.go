package main

import (
	"net/http"
	"server/db"
	"server/conn"
)


func main() {

	session := db.StartConn("localhost:27017")
	defer db.CloseConn(session)

	http.HandleFunc("/getAllMainGroups", conn.GetAll(db.Coll["main_group"]))
	http.HandleFunc("/getMainGroups", conn.Get(db.Coll["main_group"]))
	http.HandleFunc("/getOneMainGroup", conn.GetOne(db.Coll["main_group"]))
	http.HandleFunc("/setMainGroup", conn.Set(db.Coll["main_group"]))
	http.HandleFunc("/updateMainGroup", conn.Update(db.Coll["main_group"]))
	http.HandleFunc("/removeMainGroup", conn.Remove(db.Coll["main_group"]))
	http.HandleFunc("/getAllSubGroups", conn.GetAll(db.Coll["sub_group"]))
	http.HandleFunc("/getSubGroups", conn.Get(db.Coll["sub_group"]))
	http.HandleFunc("/getOneSubGroup", conn.GetOne(db.Coll["sub_group"]))
	http.HandleFunc("/setSubGroup", conn.Set(db.Coll["sub_group"]))
	http.HandleFunc("/updateSubGroup", conn.Update(db.Coll["sub_group"]))
	http.HandleFunc("/removeSubGroup", conn.Remove(db.Coll["sub_group"]))
	http.HandleFunc("/getAllCriteria", conn.GetAll(db.Coll["criterion"]))
	http.HandleFunc("/getCriteria", conn.Get(db.Coll["criterion"]))
	http.HandleFunc("/getOneCriterion", conn.GetOne(db.Coll["criterion"]))
	http.HandleFunc("/setCriterion", conn.Set(db.Coll["criterion"]))
	http.HandleFunc("/updateCriterion", conn.Update(db.Coll["criterion"]))
	http.HandleFunc("/removeCriterion", conn.Remove(db.Coll["criterion"]))
	http.HandleFunc("/getAllAccessibilities", conn.GetAll(db.Coll["accessibility"]))
	http.HandleFunc("/getAccessibilities", conn.Get(db.Coll["accessibility"]))
	http.HandleFunc("/getOneAccessibility", conn.GetOne(db.Coll["accessibility"]))
	http.HandleFunc("/setAccessibility", conn.Set(db.Coll["accessibility"]))
	http.HandleFunc("/updateAccessibility", conn.Update(db.Coll["accessibility"]))
	http.HandleFunc("/removeAccessibility", conn.Remove(db.Coll["accessibility"]))
	http.HandleFunc("/property", conn.GetHandlerProperty(db.Coll["property"]))

	http.ListenAndServe(":8080", nil)
}

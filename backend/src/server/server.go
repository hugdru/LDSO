package main

import (
	"net/http"
	"server/conn"
	"server/db"
//	"github.com/pressly/chi"
)

func main() {

	session := db.Connect("mongodb:27017")
	defer db.Disconnect(session)

	http.HandleFunc("/getAllMainGroups", conn.GetAll(db.GetCollection("main_group")))
	http.HandleFunc("/getMainGroups", conn.Get(db.GetCollection("main_group")))
	http.HandleFunc("/getOneMainGroup", conn.GetOne(db.GetCollection("main_group")))
	http.HandleFunc("/setMainGroup", conn.Set(db.GetCollection("main_group")))
	http.HandleFunc("/updateMainGroup", conn.Update(db.GetCollection("main_group")))
	http.HandleFunc("/removeMainGroup", conn.Remove(db.GetCollection("main_group")))
	http.HandleFunc("/getAllSubGroups", conn.GetAll(db.GetCollection("sub_group")))
	http.HandleFunc("/getSubGroups", conn.Get(db.GetCollection("sub_group")))
	http.HandleFunc("/getOneSubGroup", conn.GetOne(db.GetCollection("sub_group")))
	http.HandleFunc("/setSubGroup", conn.Set(db.GetCollection("sub_group")))
	http.HandleFunc("/updateSubGroup", conn.Update(db.GetCollection("sub_group")))
	http.HandleFunc("/removeSubGroup", conn.Remove(db.GetCollection("sub_group")))
	http.HandleFunc("/getAllCriteria", conn.GetAll(db.GetCollection("criterion")))
	http.HandleFunc("/getCriteria", conn.Get(db.GetCollection("criterion")))
	http.HandleFunc("/getOneCriterion", conn.GetOne(db.GetCollection("criterion")))
	http.HandleFunc("/setCriterion", conn.Set(db.GetCollection("criterion")))
	http.HandleFunc("/updateCriterion", conn.Update(db.GetCollection("criterion")))
	http.HandleFunc("/removeCriterion", conn.Remove(db.GetCollection("criterion")))
	http.HandleFunc("/getAllAccessibilities", conn.GetAll(db.GetCollection("accessibility")))
	http.HandleFunc("/getAccessibilities", conn.Get(db.GetCollection("accessibility")))
	http.HandleFunc("/getOneAccessibility", conn.GetOne(db.GetCollection("accessibility")))
	http.HandleFunc("/setAccessibility", conn.Set(db.GetCollection("accessibility")))
	http.HandleFunc("/updateAccessibility", conn.Update(db.GetCollection("accessibility")))
	http.HandleFunc("/removeAccessibility", conn.Remove(db.GetCollection("accessibility")))
	http.HandleFunc("/property", conn.GetHandlerProperty(db.GetCollection("property")))

	http.ListenAndServe(":8080", nil)
}

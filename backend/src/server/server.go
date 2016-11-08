package main

import (
	"net/http"
	"server/conn"
	"server/db"
	"github.com/pressly/chi"
)

func main() {

	dbSession := db.Connect("mongodb:27017")
	defer db.Disconnect(dbSession)

	db.SetCollection(dbSession, "Places4All", "main_group")
	db.SetCollection(dbSession, "Places4All", "sub_group")
	db.SetCollection(dbSession, "Places4All", "criterion")
	db.SetCollection(dbSession, "Places4All", "accessibility")
	db.SetCollection(dbSession, "Places4All", "property")

	router := chi.NewRouter()

	router.Get("/mainGroup", conn.Get(db.GetCollection("main_group")))
	router.Get("/mainGroup/find", conn.GetOne(db.GetCollection("main_group")))
	router.Post("/mainGroup", conn.Set(db.GetCollection("main_group")))
	router.Put("/mainGroup", conn.Update(db.GetCollection("main_group")))
	router.Delete("/mainGroup", conn.Remove(db.GetCollection("main_group")))

	router.Get("/subGroup", conn.Get(db.GetCollection("sub_group")))
	router.Get("/subGroup/find", conn.GetOne(db.GetCollection("sub_group")))
	router.Post("/subGroup", conn.Set(db.GetCollection("sub_group")))
	router.Put("/subGroup", conn.Update(db.GetCollection("sub_group")))
	router.Delete("/subGroup", conn.Remove(db.GetCollection("sub_group")))

	router.Get("/criterion", conn.Get(db.GetCollection("criterion")))
	router.Get("/criterion/find", conn.GetOne(db.GetCollection("criterion")))
	router.Post("/criterion", conn.Set(db.GetCollection("criterion")))
	router.Put("/criterion", conn.Update(db.GetCollection("criterion")))
	router.Delete("/criterion", conn.Remove(db.GetCollection("criterion")))

	router.Get("/accessibility", conn.Get(db.GetCollection("accessibility")))
	router.Get("/accessibility/find", conn.GetOne(db.GetCollection("accessibility")))
	router.Post("/accessibility", conn.Set(db.GetCollection("accessibility")))
	router.Put("/accessibility", conn.Update(db.GetCollection("accessibility")))
	router.Delete("/accessibility", conn.Remove(db.GetCollection("accessibility")))

	/*
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
	*/
	http.HandleFunc("/property", conn.GetHandlerProperty(db.GetCollection("property")))

	http.ListenAndServe(":8080", router)
}

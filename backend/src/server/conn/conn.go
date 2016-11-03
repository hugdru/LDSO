package conn

import (
	"log"
	"encoding/json"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/data"
	"server/db"
	"strconv"
)

func giveAccess(w http.ResponseWriter, methods string) {
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers",
		"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
}

func allowOrigin(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}

func GetHandlerProperty(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arg := data.Property{}

		allowOrigin(w, r)

		label := r.FormValue("label")
		value := r.FormValue("value")

		err := coll.Find(bson.M{label: value}).One(&arg)
		if err != nil {
			log.Fatal(err)
		}

		giveAccess(w, "GET, POST")

		err = json.NewEncoder(w).Encode(arg);
		if err != nil {
			log.Println(err)
		}
		log.Println(arg)
	}
}

func GetValue(r *http.Request) interface {} {
	value_type := r.FormValue("type")
	var value interface{}
	var err error
	if value_type == "int" {
		value, err = strconv.ParseInt(r.FormValue("value"), 10, 64)
		if err != nil {
			log.Panic(err)
		}
	} else {
		value = r.FormValue("value")
	}
	return value
}

func GetDocument(coll_name string) interface{} {
	var document interface{}
	switch coll_name {
	case "main_group":
		document = data.Main_Group{}
	case "sub_group":
		document = data.Sub_Group{}
	case "criterion":
		document = data.Criterion{}
	case "accessibility":
		document = data.Accessibility{}
	}
	return document
}

func GetDocuments(coll *mgo.Collection, tagged bool, tag string, value interface{}) interface{} {
	var document interface{}
	switch coll.Name {
	case "main_group":
		main_groups := []data.Main_Group{}
		db.Find(coll, &main_groups, tagged, tag, value)
		document = main_groups
	case "sub_group":
		sub_groups := []data.Sub_Group{}
		db.Find(coll, &sub_groups, tagged, tag, value)
		document = sub_groups
	case "criterion":
		criteria := []data.Criterion{}
		db.Find(coll, &criteria, tagged, tag, value)
		document = criteria
	case "accessibility":
		accessibilities := []data.Accessibility{}
		db.Find(coll, &accessibilities, tagged, tag, value)
		document = accessibilities
	}
	return document
}

func GetAll(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documents := GetDocuments(coll, false, "", 0)
		switch coll.Name {
		}
		giveAccess(w, "GET, POST")
		err := json.NewEncoder(w).Encode(documents);
		if err != nil {
			log.Panic(err)
		}
	}
}

func Get(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.FormValue("tag")
		value := GetValue(r)
		documents := GetDocuments(coll, true, tag, value)
		giveAccess(w, "GET, POST")
		err := json.NewEncoder(w).Encode(documents);
		if err != nil {
			log.Panic(err)
		}
	}
}

func GetOne(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document := GetDocument(coll.Name)
		tag := r.FormValue("tag")
		value := GetValue(r)
		db.FindOne(coll, &document, tag, value)
		giveAccess(w, "GET, POST")
		err := json.NewEncoder(w).Encode(document);
		if err != nil {
			log.Panic(err)
		}
	}
}
func Set(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document := GetDocument(coll.Name)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&document)
		if err != nil {
			log.Panic(err)
		}
		db.Insert(coll, &document)
	}
}

func Update(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document := GetDocument(coll.Name)
		id, err := strconv.ParseInt(r.FormValue("_id"), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		db.FindOne(coll, &document, "_id", id)
		tag := r.FormValue("tag")
		value := GetValue(r)
		db.Update(coll, document, tag, value)
	}
}

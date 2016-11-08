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
		value, err = strconv.Atoi(r.FormValue("value"))
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

func Get(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowOrigin(w, r)
		tag := r.FormValue("tag")
		if tag == "" {
			documents := GetDocuments(coll, false, "", 0)
			err := json.NewEncoder(w).Encode(documents);
			if err != nil {
				log.Panic(err)
			}
			log.Println(documents)
		} else {
			value := GetValue(r)
			documents := GetDocuments(coll, true, tag, value)
			err := json.NewEncoder(w).Encode(documents);
			if err != nil {
				log.Panic(err)
			}
			log.Println(documents)
		}
	}
}

func GetOne(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowOrigin(w, r)
		document := GetDocument(coll.Name)
		tag := r.FormValue("tag")
		value := GetValue(r)
		db.FindOne(coll, &document, tag, value)
		err := json.NewEncoder(w).Encode(document);
		if err != nil {
			log.Panic(err)
		}
		log.Println(document)
	}
}

func Set(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowOrigin(w, r)
		document := GetDocument(coll.Name)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&document)
		if err != nil {
			log.Panic(err)
		}
		db.Insert(coll, &document)
		log.Println(document)
	}
}

func Update(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowOrigin(w, r)
		document := GetDocument(coll.Name)
		id, err := strconv.Atoi(r.FormValue("_id"))
		if err != nil {
			log.Panic(err)
		}
		db.FindOne(coll, &document, "_id", id)
		tag := r.FormValue("tag")
		value := GetValue(r)
		db.Update(coll, document, tag, value)
		db.FindOne(coll, &document, "_id", id)
		log.Println(document)
	}
}

func RecursiveRemove(coll *mgo.Collection, id int) {
	document := GetDocument(coll.Name)
	db.FindOne(coll, &document, "_id", id)
	db.Remove(coll, "_id", id)
	log.Println(document)
	var child_coll *mgo.Collection
	switch coll.Name {
	case "main_group":
		sub_groups := []data.Sub_Group{}
		child_coll = db.GetCollection("sub_group");
		db.Find(child_coll, &sub_groups, true, "main_group", id)
		for _, sub_group := range sub_groups {
			RecursiveRemove(child_coll, sub_group.Id)
		}
	case "sub_group":
		criteria := []data.Criterion{}
		child_coll = db.GetCollection("criterion")
		db.Find(child_coll, &criteria, true, "sub_group", id)
		for _, criterion := range criteria {
			RecursiveRemove(child_coll, criterion.Id)
		}
	case "criterion":
		accessibilities := []data.Accessibility{}
		child_coll = db.GetCollection("accessibility")
		db.Find(child_coll, &accessibilities, true, "criterion", id)
		for _, accessibility := range accessibilities {
			RecursiveRemove(child_coll, accessibility.Id)
		}
	}
}

func Remove(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowOrigin(w, r)
		id, err := strconv.Atoi(r.FormValue("_id"))
		if err != nil {
			log.Panic(err)
		}
		RecursiveRemove(coll, id)
	}
}

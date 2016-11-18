package conn

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"server/data"
	"server/db"
	"strconv"
)

func Decode(r *http.Request, document interface{}) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&document)
	if err != nil {
		log.Panic(err)
	}
}

func GetValue(r *http.Request) interface{} {
	var value interface{}
	var err error

	value_type := r.FormValue("type")
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
	case "property":
		document = data.Property{}
	case "audit":
		document = data.Audit{}
	}
	return document
}

func GetDocuments(coll *mgo.Collection, tagged bool, tag string,
	value interface{}) interface{} {
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
	case "property":
		properties := []data.Property{}
		db.Find(coll, &properties, tagged, tag, value)
		document = properties
	case "audit":
		audits := []data.Audit{}
		db.Find(coll, &audits, tagged, tag, value)
		document = audits
	}

	return document
}

func GetIds(coll *mgo.Collection) []int {
	var ids []int
	switch coll.Name {
	case "main_group":
		main_groups := []data.Main_Group{}
		db.Find(coll, &main_groups, false, "", 0)
		for _, main_group := range main_groups {
			ids = append(ids, main_group.Id)
		}
	case "sub_group":
		sub_groups := []data.Sub_Group{}
		db.Find(coll, &sub_groups, false, "", 0)
		for _, sub_group := range sub_groups {
			ids = append(ids, sub_group.Id)
		}
	case "criterion":
		criteria := []data.Criterion{}
		db.Find(coll, &criteria, false, "", 0)
		for _, criterion := range criteria {
			ids = append(ids, criterion.Id)
		}
	case "accessibility":
		accessibilities := []data.Accessibility{}
		db.Find(coll, &accessibilities, false, "", 0)
		for _, accessibility := range accessibilities {
			ids = append(ids, accessibility.Id)
		}
	case "property":
		properties := []data.Property{}
		db.Find(coll, &properties, false, "", 0)
		for _, property := range properties {
			ids = append(ids, property.Id)
		}
	case "audit":
		audits := []data.Audit{}
		db.Find(coll, &audits, false, "", 0)
		for _, audit := range audits {
			ids = append(ids, audit.Id)
		}
	}
	return ids
}

func Get(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.FormValue("tag")
		if tag == "" {
			documents := GetDocuments(coll, false, "", 0)
			err := json.NewEncoder(w).Encode(documents)
			if err != nil {
				log.Panic(err)
			}
			log.Println(documents)
		} else {
			value := GetValue(r)
			documents := GetDocuments(coll, true, tag, value)
			err := json.NewEncoder(w).Encode(documents)
			if err != nil {
				log.Panic(err)
			}

			log.Println(documents)
		}
	}
}

func GetOne(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document := GetDocument(coll.Name)
		tag := r.FormValue("tag")
		value := GetValue(r)
		db.FindOne(coll, &document, tag, value)
		err := json.NewEncoder(w).Encode(document)
		if err != nil {
			log.Panic(err)
		}
		log.Println(document)
	}
}

func Set(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var err error
		defer r.Body.Close()
		var id int
		switch coll.Name {
		case "main_group":
			documentMaxId := data.Main_Group{}
			document := data.Main_Group{}
			db.FindMaxId(coll, &documentMaxId)
			err = decoder.Decode(&document)
			if err != nil {
				log.Panic(err)
			}
			id = documentMaxId.Id + 1
			document.Id = id
			db.Insert(coll, document)
			log.Println(document)
		case "sub_group":
			documentMaxId := data.Sub_Group{}
			document := data.Sub_Group{}
			db.FindMaxId(coll, &documentMaxId)
			err = decoder.Decode(&document)
			if err != nil {
				log.Panic(err)
			}
			id = documentMaxId.Id + 1
			document.Id = id
			db.Insert(coll, document)
			log.Println(document)
		case "criterion":
			documentMaxId := data.Criterion{}
			document := data.Criterion{}
			db.FindMaxId(coll, &documentMaxId)
			err = decoder.Decode(&document)
			if err != nil {
				log.Panic(err)
			}
			id = documentMaxId.Id + 1
			document.Id = id
			db.Insert(coll, document)
			log.Println(document)
		case "accessibility":
			documentMaxId := data.Accessibility{}
			document := data.Accessibility{}
			db.FindMaxId(coll, &documentMaxId)
			err = decoder.Decode(&document)
			if err != nil {
				log.Panic(err)
			}
			id = documentMaxId.Id + 1
			document.Id = id
			db.Insert(coll, document)
			log.Println(document)
		case "property":
			documentMaxId := data.Property{}
			document := data.Property{}
			db.FindMaxId(coll, &documentMaxId)
			err = decoder.Decode(&document)
			if err != nil {
				log.Panic(err)
			}
			id = documentMaxId.Id + 1
			document.Id = id
			db.Insert(coll, document)
			log.Println(document)
		case "audit":
			documentMaxId := data.Audit{}
			document := data.Audit{}
			db.FindMaxId(coll, &documentMaxId)
			err = decoder.Decode(&document)
			if err != nil {
				log.Panic(err)
			}
			id = documentMaxId.Id + 1
			document.Id = id
			db.Insert(coll, document)
			log.Println(document)
		}
		err = json.NewEncoder(w).Encode(id)
		if err != nil {
			log.Panic(err)
		}

	}
}

func Update(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newDocument := GetDocument(coll.Name)
		oldDocument := GetDocument(coll.Name)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&newDocument)
		if err != nil {
			log.Panic(err)
		}
		id, err := strconv.Atoi(r.FormValue("_id"))
		if err != nil {
			log.Panic(err)
		}
		db.FindOne(coll, &oldDocument, "_id", id)
		db.Update(coll, oldDocument, newDocument)
		db.FindOne(coll, &oldDocument, "_id", id)
		log.Println(oldDocument)
	}
}

func RecursiveDelete(coll *mgo.Collection, id int) {
	var child_coll *mgo.Collection
	document := GetDocument(coll.Name)
	db.FindOne(coll, &document, "_id", id)
	db.Delete(coll, "_id", id)
	log.Println(document)
	switch coll.Name {
	case "main_group":
		sub_groups := []data.Sub_Group{}
		child_coll = db.GetCollection("sub_group")
		db.Find(child_coll, &sub_groups, true, "main_group", id)
		for _, sub_group := range sub_groups {
			RecursiveDelete(child_coll, sub_group.Id)
		}
	case "sub_group":
		criteria := []data.Criterion{}
		child_coll = db.GetCollection("criterion")
		db.Find(child_coll, &criteria, true, "sub_group", id)
		for _, criterion := range criteria {
			RecursiveDelete(child_coll, criterion.Id)
		}
	case "criterion":
		accessibilities := []data.Accessibility{}
		child_coll = db.GetCollection("accessibility")
		db.Find(child_coll, &accessibilities, true, "criterion", id)
		for _, accessibility := range accessibilities {
			RecursiveDelete(child_coll, accessibility.Id)
		}
	}
}

func Delete(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id_str := r.FormValue("_id")
		if id_str == "" {
			ids := GetIds(coll)
			for _, id := range ids {
				RecursiveDelete(coll, id)
			}
		} else {
			id, err := strconv.Atoi(id_str)
			if err != nil {
				log.Panic(err)
			}
			RecursiveDelete(coll, id)
		}
	}
}

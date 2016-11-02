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

func GetAll(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var document interface{}
		switch coll.Name {
		case "main_group":
			main_group := []data.Main_Group{}
			db.FindAll(coll, &main_group)
			document = main_group
		case "sub_group":
			sub_group := []data.Sub_Group{}
			db.FindAll(coll, &sub_group)
			document = sub_group
		case "criterion":
			criterion := []data.Criterion{}
			db.FindAll(coll, &criterion)
			document = criterion
		case "accessibility":
			accessibility := []data.Accessibility{}
			db.FindAll(coll, &accessibility)
			document = accessibility
		}
		giveAccess(w, "GET, POST")
		err := json.NewEncoder(w).Encode(document);
		if err != nil {
			log.Panic(err)
		}
	}
}

func Get(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.FormValue("tag")
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
		var document interface{}
		switch coll.Name {
		case "main_group":
			main_groups := []data.Main_Group{}
			db.Find(coll, &main_groups, tag, value)
			document = main_groups
		case "sub_group":
			sub_groups := []data.Sub_Group{}
			db.Find(coll, &sub_groups, tag, value)
			document = sub_groups
		case "criterion":
			criteria := []data.Criterion{}
			db.Find(coll, &criteria, tag, value)
			document = criteria
		case "accessibility":
			accessibilities := []data.Accessibility{}
			db.Find(coll, &accessibilities, tag, value)
			document = accessibilities
		}
		giveAccess(w, "GET, POST")
		err = json.NewEncoder(w).Encode(document);
		if err != nil {
			log.Panic(err)
		}
	}
}

func GetOne(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document := GetDocument(coll.Name)
		tag := r.FormValue("tag")
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
		db.FindOne(coll, &document, tag, value)
		giveAccess(w, "GET, POST")
		err = json.NewEncoder(w).Encode(document);
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
		value_type := r.FormValue("type")
		var value interface{}
		if value_type == "int" {
			value, err = strconv.ParseInt(r.FormValue("value"), 10, 64)
			if err != nil {
				log.Panic(err)
			}
		} else {
			value = r.FormValue("value")
		}
		db.Update(coll, document, tag, value)
	}
}

/*
func SetSubGroup(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub_group := data.Sub_Group{}

		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&sub_group)
		if err != nil {
			log.Panic(err)
		}

		main_group := data.Main_Group{}
		main_group_name:= r.FormValue("main_group")
		db.FindOne(coll, &main_group, "name", main_group_name)

		db.Add(coll, &main_group, "sub_groups", &sub_group)
	}
}


type SubgroupContainer struct {
	Sub_Groups []data.Sub_Group
}

func SetCriterion(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		criterion := data.Criterion{}

		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&criterion)
		if err != nil {
			log.Panic(err)
		}

		main_group_name := r.FormValue("main_group")
		sub_group_name := r.FormValue("sub_group")

		container := SubgroupContainer{}
		err = coll.Find(bson.M{"name": main_group_name}).Select(bson.M{"sub_groups": bson.M{"$elemMatch": bson.M{"name": sub_group_name}}}).One(&container)
		if err != nil {
			log.Panic(err)
		}
		sub_group := container.Sub_Groups[0]
		log.Println(sub_group)

		change := bson.M{"$addToSet": bson.M{"sub_groups.$.criteria": &criterion}}
		log.Println(criterion)
		log.Println(change)
		err = coll.Update(&sub_group, change)
		if err != nil {
			log.Panic(err)
		}

		main_group := data.Main_Group{}
		db.FindOne(coll, &main_group, "name", main_group_name)
		log.Println(main_group)
	}
}
*/
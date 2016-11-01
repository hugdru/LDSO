package conn

import (
	"log"
	"encoding/json"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/data"
	"server/db"
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
	return func (w http.ResponseWriter, r *http.Request) {
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

func GetAllGroups(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		groups := []data.Main_Group{}

//		allowOrigin(w, r)

		db.FindAll(coll, &groups)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(groups);
		if err != nil {
			log.Println(err)
		}
		log.Println(groups)
	}
}

type name_query struct {
	Name string
}

func GetGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		group := data.Main_Group{}
		name := name_query{}

//		allowOrigin(w, r)

		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&name)
		if err != nil {
			log.Panic(err)
		}


		log.Println(r)
		log.Println(name.Name)

		db.FindOne(coll, &group, "name", name.Name)

		giveAccess(w, "GET, POST")

		err = json.NewEncoder(w).Encode(group);
		if err != nil {
			log.Println(err)
		}
		log.Println(group)
	}
}

func SetGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var group data.Main_Group
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&group)
		if err != nil {
			log.Panic(err)
		}

		db.Insert(coll, &group)

		log.Println(group)
	}
}

func SetSubGroup(coll *mgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var subgroup data.Sub_Group
		var group data.Main_Group
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&subgroup)
		if err != nil {
			log.Panic(err)
		}

		log.Println(group)
	}
}
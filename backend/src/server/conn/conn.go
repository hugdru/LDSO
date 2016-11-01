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

/*		name := name_query{}

		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&name)
		if err != nil {
			log.Panic(err)
		}
*/

		name:= r.FormValue("name")

		db.FindOne(coll, &group, "name", name)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(group);
		if err != nil {
			log.Println(err)
		}
		log.Println(group)
	}
}

func SetGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		group := data.Main_Group{}
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
		sub_group := data.Sub_Group{}

		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&sub_group)
		if err != nil {
			log.Panic(err)
		}

		group := data.Main_Group{}

		name:= r.FormValue("name")

		db.FindOne(coll, &group, "name", name)

		db.Add(coll, &group, "sub_groups", &sub_group)

		db.FindOne(coll, &group, "name", name)
		log.Println(group)
	}
}
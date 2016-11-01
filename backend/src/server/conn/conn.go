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
		arg := []data.Group{}

		allowOrigin(w, r)

		err := coll.Find(bson.M{}).Iter().All(&arg)
		if err != nil {
			log.Panic(err)
		}

		giveAccess(w, "GET, POST")

		err = json.NewEncoder(w).Encode(arg);
		if err != nil {
			log.Println(err)
		}
		log.Println(arg)
	}
}

func SetGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var group data.Group
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

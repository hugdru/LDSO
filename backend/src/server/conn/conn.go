package conn

import (
	"log"
	"encoding/json"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/data"
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

func GetHandler(arg interface{}) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		allowOrigin(w, r)
		giveAccess(w, "GET")

		err := json.NewEncoder(w).Encode(arg);
		if err != nil {
			panic(err)
		}
		log.Println(arg)
	}
}

func GetHandlerProperty(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		arg := data.Property{}

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		label := r.FormValue("label")
		value := r.FormValue("value")

		err := coll.Find(bson.M{label: value}).One(&arg)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

		err = json.NewEncoder(w).Encode(arg);
		if err != nil {
			panic(err)
		}
		log.Println(arg)
	}
}

func GetHandlerCriterion(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		arg := []data.Criterion{}

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		err := coll.Find(bson.M{}).Iter().All(&arg)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

		err = json.NewEncoder(w).Encode(arg);
		if err != nil {
			panic(err)
		}
		log.Println(arg)
	}
}

func GetHandlerCriteriaSet(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		arg := data.Criterion_Set{}

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		label := r.FormValue("label")
		value := r.FormValue("value")

		err := coll.Find(bson.M{label: value}).One(&arg)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

		err = json.NewEncoder(w).Encode(arg);
		if err != nil {
			panic(err)
		}
		log.Println(arg)
	}
}


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

func GetAllMainGroups(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		main_groups := []data.Main_Group{}

		db.FindAll(coll, &main_groups)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(main_groups);
		if err != nil {
			log.Println(err)
		}
		log.Println(main_groups)
	}
}

func GetMainGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		main_group := data.Main_Group{}

		name:= r.FormValue("name")

		db.FindOne(coll, &main_group, "name", name)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(main_group);
		if err != nil {
			log.Println(err)
		}
		log.Println(main_group)
	}
}

func SetMainGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		main_group := data.Main_Group{}
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&main_group)
		if err != nil {
			log.Panic(err)
		}

		db.Insert(coll, &main_group)

		log.Println(main_group)
	}
}

func GetSubGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		sub_group := data.Sub_Group{}

		name:= r.FormValue("name")

		db.FindOne(coll, &sub_group, "name", name)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(sub_group);
		if err != nil {
			log.Println(err)
		}
		log.Println(sub_group)
	}
}

func SetSubGroup(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		sub_group := data.Sub_Group{}
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&sub_group)
		if err != nil {
			log.Panic(err)
		}

		db.Insert(coll, &sub_group)

		log.Println(sub_group)
	}
}

func GetCriterion(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		criterion := data.Criterion{}

		name:= r.FormValue("name")

		db.FindOne(coll, &criterion, "name", name)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(criterion);
		if err != nil {
			log.Println(err)
		}
		log.Println(criterion)
	}
}

func SetCriterion(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		criterion := data.Criterion{}
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&criterion)
		if err != nil {
			log.Panic(err)
		}

		db.Insert(coll, &criterion)

		log.Println(criterion)
	}
}

func GetAccessibility(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		accessibility := data.Accessibility{}

		name:= r.FormValue("name")

		db.FindOne(coll, &accessibility, "name", name)

		giveAccess(w, "GET, POST")

		err := json.NewEncoder(w).Encode(accessibility);
		if err != nil {
			log.Println(err)
		}
		log.Println(accessibility)
	}
}

func SetAccessibility(coll *mgo.Collection) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		accessibility := data.Accessibility{}
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		err := decoder.Decode(&accessibility)
		if err != nil {
			log.Panic(err)
		}

		db.Insert(coll, &accessibility)

		log.Println(accessibility)
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
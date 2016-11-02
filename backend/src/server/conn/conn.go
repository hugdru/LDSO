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

func GetAll(document interface{}) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		giveAccess(w, "GET, POST")
		err := json.NewEncoder(w).Encode(document);
		if err != nil {
			log.Println(err)
		}
	}
}

func GetOne(coll *mgo.Collection, document interface{}) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		tag := r.FormValue("tag")
		value_type := r.FormValue("type")
		if value_type == "int" {
			value, err := strconv.ParseInt(r.FormValue("value"), 10, 64)
			if err != nil {
				log.Panic(err)
			}
			db.FindOne(coll, &document, tag, value)
		} else {
			value := r.FormValue("value")
			db.FindOne(coll, &document, tag, value)
		}
		giveAccess(w, "GET, POST")
		err := json.NewEncoder(w).Encode(document);
		if err != nil {
			log.Println(err)
		}
	}
}

func Set(coll *mgo.Collection, document interface{}) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&document)
		if err != nil {
			log.Panic(err)
		}
		db.Insert(coll, &document)
	}
}

func Update(coll *mgo.Collection, document interface{}) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.FormValue("_id"), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		db.FindOne(coll, &document, "_id", id)
		tag := r.FormValue("tag")
		value_type := r.FormValue("type")
		if value_type == "int" {
			value, err := strconv.ParseInt(r.FormValue("value"), 10, 64)
			if err != nil {
				log.Panic(err)
			}
			db.Update(coll, document, tag, value)
		} else {
			value := r.FormValue("value")
			db.Update(coll, document, tag, value)
		}
	}
}

func GetAllMainGroups(coll *mgo.Collection) http.HandlerFunc {
	main_groups := []data.Main_Group{}
	db.FindAll(coll, &main_groups)
	return GetAll(main_groups)
}

func GetMainGroup(coll *mgo.Collection) http.HandlerFunc {
	main_group := data.Main_Group{}
	return GetOne(coll, main_group)
}

func SetMainGroup(coll *mgo.Collection) http.HandlerFunc {
	main_group := data.Main_Group{}
	return Set(coll, main_group)
}

func UpdateMainGroup(coll *mgo.Collection) http.HandlerFunc {
	main_group := data.Main_Group{}
	return Update(coll, main_group)
}

func GetAllSubGroups(coll *mgo.Collection) http.HandlerFunc {
	sub_groups := []data.Sub_Group{}
	db.FindAll(coll, &sub_groups)
	return GetAll(sub_groups)
}

func GetSubGroup(coll *mgo.Collection) http.HandlerFunc {
	sub_group := data.Sub_Group{}
	return GetOne(coll, sub_group)
}

func SetSubGroup(coll *mgo.Collection) http.HandlerFunc {
	sub_group := data.Sub_Group{}
	return Set(coll, sub_group)
}

func UpdateSubGroup(coll *mgo.Collection) http.HandlerFunc {
	sub_group := data.Sub_Group{}
	return Update(coll, sub_group)
}

func GetAllCriteria(coll *mgo.Collection) http.HandlerFunc {
	criteria := []data.Criterion{}
	db.FindAll(coll, &criteria)
	return GetAll(criteria)
}

func GetCriterion(coll *mgo.Collection) http.HandlerFunc {
	criterion := data.Criterion{}
	return GetOne(coll, criterion)
}

func SetCriterion(coll *mgo.Collection) http.HandlerFunc {
	criterion := data.Criterion{}
	return Set(coll, criterion)
}

func UpdateCriterion(coll *mgo.Collection) http.HandlerFunc {
	criterion := data.Criterion{}
	return Update(coll, criterion)
}

func GetAllAccessibilities(coll *mgo.Collection) http.HandlerFunc {
	accessibilities := []data.Accessibility{}
	db.FindAll(coll, &accessibilities)
	return GetAll(accessibilities)
}

func GetAccessibility(coll *mgo.Collection) http.HandlerFunc {
	accessibility := data.Accessibility{}
	return GetOne(coll, accessibility)
}

func SetAccessibility(coll *mgo.Collection) http.HandlerFunc {
	accessibility := data.Accessibility{}
	return Set(coll, accessibility)
}

func UpdateAccessibility(coll *mgo.Collection) http.HandlerFunc {
	accessibility := data.Accessibility{}
	return Update(coll, accessibility)
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
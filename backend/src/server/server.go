package main

import (
	"net/http"
	"server/db"
	"server/conn"
	"server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	// TODO: Find a better way to make the index unique on first run.
	// Should be done in a script not on go.
	first_run := false

	session := db.StartConn("localhost:27017")
	defer db.CloseConn(session)

	if !db.ExistsCollections(session, "Places4All") {
		first_run = true
	}

	c_criterion := db.GetCollection(session, "Places4All", "Criterion")
	c_sub_criterion := db.GetCollection(session, "Places4All", "Sub_Criterion")
	c_property := db.GetCollection(session, "Places4All", "Property")
	c_criterion_set := db.GetCollection(session, "Places4All", "Criterion_Set")
	// c_evaluation := db.GetCollection(session, "Places4All", "Evaluation")
	// c_note := db.GetCollection(session, "Places4All", "Note")

	if first_run == true {
		db.EnsureUnique(c_criterion, "name")
		db.EnsureUnique(c_sub_criterion, "name")
		db.EnsureUnique(c_property, "name")
	}

	populateCriterion(c_criterion);
	populateSubCriterion(c_sub_criterion);
	populateProperty(c_property);
	populateCriterionSet(c_criterion_set, c_criterion, c_sub_criterion)

	http.HandleFunc("/criterion", conn.GetHandlerCriterion(c_criterion))
	// http.HandleFunc("/property", conn.GetHandlerProperty(c_property))

	http.ListenAndServe(":8080", nil)
}

/*
 * Test data start
 * Delete latter
 */
func populateSubCriterion(c_sub_criterion *mgo.Collection) {
	sub1 := data.Sub_Criterion{"Rampa", 20}
	sub2 := data.Sub_Criterion{"Porta", 20}
	sub3 := data.Sub_Criterion{"Altura", 60}

	sub4 := data.Sub_Criterion{"Entrada", 10}
	sub5 := data.Sub_Criterion{"WC", 40}
	sub6 := data.Sub_Criterion{"Cadeira", 50}

	db.Insert(c_sub_criterion, &sub1, &sub2, &sub3, &sub4, &sub5, &sub6)
}

func populateCriterion(c_criterion *mgo.Collection) {
	crit1 := data.Criterion{"Acessos", 40}
	crit2 := data.Criterion{"Percurso Exterior", 20}
	crit3 := data.Criterion{"Percurso Interior", 20}
	crit4 := data.Criterion{"Bens e Serviços", 20}

	db.Insert(c_criterion, &crit1, &crit2, &crit3, &crit4)
}

func populateProperty(c_property *mgo.Collection) {
	prop1 := data.Property{"Casa", data.Owner{"Joao"}, "pic1.png"}
	prop2 := data.Property{"Hotel Sunny", data.Owner{"Carlos"}, "pic2.png"}

	db.Insert(c_property, &prop1, &prop2)
}

func populateCriterionSet(c_crit_set, c_crit, c_sub *mgo.Collection) {
	var crit data.Criterion
	var sub1, sub2, sub3 data.Sub_Criterion
	var crit_set1 data.Criterion_Set

	c_crit.Find(bson.M{"name": "Acessos"}).One(&crit)
	c_sub.Find(bson.M{"name": "Entrada"}).One(&sub1)
	c_sub.Find(bson.M{"name": "WC"}).One(&sub2)
	c_sub.Find(bson.M{"name": "Cadeira"}).One(&sub3)

	crit_set1.Criterion = crit
	crit_set1.SetSub_Criterion(sub1, sub2, sub3)

	db.Insert(c_crit_set, &crit_set1)
}
/*
 * Test data end
 */

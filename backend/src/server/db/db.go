package db

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	Coll map[string]*mgo.Collection
)

type InsertFunc func(c *mgo.Collection, args ...interface{})

func StartConn(addr string) *mgo.Session {
	session, err := mgo.Dial(addr)
	if err != nil {
		panic(err)
	}
	Coll = make(map[string]*mgo.Collection)
	Coll["main_group"] = GetCollection(session, "Places4All", "main_group")
	Coll["sub_group"] = GetCollection(session, "Places4All", "sub_group")
	Coll["criterion"] = GetCollection(session, "Places4All", "criterion")
	Coll["accessibility"] = GetCollection(session, "Places4All", "accessibility")
	Coll["property"] = GetCollection(session, "Places4All", "property")
	return session
}

func CloseConn(session *mgo.Session) {
	session.Close()
}

func GetCollection(session *mgo.Session, db_name,
		c_name string) *mgo.Collection {
	return session.DB(db_name).C(c_name)
}

func Insert(c *mgo.Collection, documents ...interface{}) {
	err := c.Insert(documents...)
	if err != nil {
		log.Panic(err)
	}
}

func FindOne(c *mgo.Collection, document interface{}, tag string, value interface{}) {
	err := c.Find(bson.M{tag: value}).One(document)
	if err != nil {
		log.Panic(err)
	}
}

func Find(c *mgo.Collection, documents interface{}, tagged bool, tag string, value interface{}) {
	var err error
	if (tagged) {
		err = c.Find(bson.M{tag: value}).Iter().All(documents)
	} else {
		err = c.Find(bson.M{}).Iter().All(documents)
	}
	if err != nil {
		log.Panic(err)
	}
}

func Update(c *mgo.Collection, document interface{}, tag string, value interface{}) {
	change := bson.M{"$set": bson.M{tag: value}}
	err := c.Update(document, change)
	if err != nil {
		log.Panic(err)
	}
}

func Remove(c *mgo.Collection, tag string, value interface{}) {
	err := c.Remove(bson.M{tag: value})
	if err != nil {
		log.Panic(err)

	}
}

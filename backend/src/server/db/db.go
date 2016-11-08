package db

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	coll map[string]*mgo.Collection
)

type InsertFunc func(c *mgo.Collection, args ...interface{})

func Connect(addr string) *mgo.Session {
	session, err := mgo.Dial(addr)
	if err != nil {
		panic(err)
	}
	coll = make(map[string]*mgo.Collection)
	return session
}

func Disconnect(session *mgo.Session) {
	session.Close()
}

func SetCollection(session *mgo.Session, db_name, c_name string) {
	coll[c_name] = session.DB(db_name).C(c_name)
}

func GetCollection(c_name string) *mgo.Collection {
	return coll[c_name]
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

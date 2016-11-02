package db

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type InsertFunc func(c *mgo.Collection, args ...interface{})

func StartConn(addr string) *mgo.Session {
	session, err := mgo.Dial(addr)
	if err != nil {
		panic(err)
	}
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

func Find(c *mgo.Collection, document interface{}, tag string, value interface{}) {
	err := c.Find(bson.M{tag: value}).Iter().All(document)
	if err != nil {
		log.Panic(err)
	}
}

func FindAll(c *mgo.Collection, document interface{}) {
	err := c.Find(bson.M{}).Iter().All(document)
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

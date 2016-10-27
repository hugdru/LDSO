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

func Insert(c *mgo.Collection, args ...interface{}) {
	err := c.Insert(args...)
	if err != nil {
		log.Print(err)
	}
}

func FindOne(c *mgo.Collection, arg interface{}, tag, value string) {
	err := c.Find(bson.M{tag: value}).One(&arg)
	if err != nil {
		log.Fatal(err)
	}
}

func Find(c *mgo.Collection, arg interface{}, tag, value string) {
	err := c.Find(bson.M{tag: value}).Limit(100).Iter().All(&arg)
	if err != nil {
		log.Fatal(err)
	}
}

func FindAll(c *mgo.Collection, arg []interface{}) {
	err := c.Find(bson.M{}).Limit(100).Iter().All(&arg)
	if err != nil {
		log.Fatal(err)
	}
}

func EnsureUnique(c *mgo.Collection, tag string) {
	index := mgo.Index{
		Key: []string{tag},
		Unique: true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		log.Print(err)
	}
}

func ExistsCollections(session *mgo.Session, name string) bool {
	names, err := session.DB(name).CollectionNames()
	if err != nil {
		log.Fatal(err)
	}
	return len(names) != 0
}

package datastore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Datastore struct {
	postgres *sqlx.DB
}

const connectionDetails = "host=postgres user=admin password=admin dbname=places4all sslmode=disable"

func Connect() *Datastore {
	postgres := sqlx.MustConnect("postgres", connectionDetails)
	return &Datastore{postgres}
}

func (datastore *Datastore) Close() {
	err := datastore.postgres.Close()
	if err != nil {
		panic(err)
	}
	datastore.postgres = nil
}

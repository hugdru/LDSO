package datastore

import (
	"database/sql"
	"github.com/alexedwards/scs/engine/memstore"
	"github.com/alexedwards/scs/session"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

type Datastore struct {
	postgres       *sqlx.DB
	SessionManager func(h http.Handler) http.Handler
}

const connectionDetails = "host=postgres user=admin password=admin dbname=places4all sslmode=disable"

func Connect() *Datastore {
	postgres := sqlx.MustConnect("postgres", connectionDetails)

	engine := memstore.New(10 * time.Minute)
	sessionManager := session.Manage(engine,
		//session.Domain("place4all.com"),
		session.HttpOnly(true),
		//session.Path("/"),
		session.IdleTimeout(7*time.Hour*24),
		session.Lifetime(14*time.Hour*24),
		session.Persist(true),
		//session.Secure(true),
		session.ErrorFunc(ServerError),
	)

	return &Datastore{postgres: postgres, SessionManager: sessionManager}
}

func (datastore *Datastore) Close() {
	err := datastore.postgres.Close()
	if err != nil {
		panic(err)
	}
	datastore.postgres = nil
}

func (datastore *Datastore) BeginTransaction() (*sql.Tx, error) {
	return datastore.postgres.Begin()
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Failed to connect to redis", 500)
}

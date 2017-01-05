package datastore

import (
	"database/sql"
	"github.com/alexedwards/scs/engine/redisstore"
	"github.com/alexedwards/scs/session"
	"github.com/garyburd/redigo/redis"
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

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	engine := redisstore.New(pool)
	sessionManager := session.Manage(engine,
		//session.Domain("example.org"),
		session.HttpOnly(true),
		session.Path("/"),
		session.IdleTimeout(15 * time.Hour * 24),
		session.Lifetime(time.Hour * 24),
		session.Persist(true),
		session.Secure(false),
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

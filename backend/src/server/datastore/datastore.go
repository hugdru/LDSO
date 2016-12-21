package datastore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/garyburd/redigo/redis"
	"github.com/alexedwards/scs/engine/redisstore"
	"github.com/alexedwards/scs/session"
	"net/http"
	"server/handler/helpers"
)

type Datastore struct {
	postgres *sqlx.DB
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
		//session.Domain("example.org"), // Domain is not set by default.
		session.HttpOnly(true),        	 // HttpOnly attribute is true by default.
		session.Path("/"),       	 // Path is set to "/" by default.
		session.Secure(false),           // Secure attribute is false by default.)
		session.ErrorFunc(ServerError),	 // Custom error handler
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

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, helpers.Error("Sorry, the application encountered an error"), 500)
}

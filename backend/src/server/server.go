package main

import (
	"github.com/pressly/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request handled by go."))
	})
	http.ListenAndServe(":8080", r)
}

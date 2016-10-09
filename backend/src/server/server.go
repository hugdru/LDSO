package main

import (
	"net/http"
	"os"
	"path/filepath"
	"github.com/pressly/chi"
)

func main() {
	router := chi.NewRouter()

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filesDir := filepath.Join(workDir, "../frontend/dist")
	router.FileServer("/", http.Dir(filesDir))

	http.ListenAndServe(":8080", router)
}
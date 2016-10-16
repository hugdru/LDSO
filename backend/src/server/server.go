package main

import (
	"github.com/pressly/chi"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := chi.NewRouter()

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filesDir := filepath.Join(workDir, "../frontend/dist")
	router.FileServer("/", http.Dir(filesDir))

	http.ListenAndServe(":8888", router)
}

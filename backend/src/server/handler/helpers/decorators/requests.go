package decorators

import (
	"net/http"
	"server/handler/helpers"
)

func RequestJson(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := helpers.GetContentType(r.Header.Get("Content-type"))
		if contentType != "application/json" {
			http.Error(w, helpers.Error("Expected, Content-type: application/json"), 415)
			return
		}
		f(w, r)
	}
}

func RequestMultipart(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := helpers.GetContentType(r.Header.Get("Content-type"))
		if contentType != "multipart/form-data" {
			http.Error(w, helpers.Error("Expected, Content-type: multipart/form-data"), 415)
			return
		}
		f(w, r)
	}
}

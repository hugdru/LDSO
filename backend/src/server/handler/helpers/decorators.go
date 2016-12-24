package helpers

import "net/http"

func RequestJson(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "Expected, Content-type: application/json", 415)
			return
		}
		f(w, r)
	}
}

func RequestMultipart(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-type")
		if contentType != "multipart/form-data" {
			http.Error(w, "Expected, Content-type: multipart/form-data", 415)
			return
		}
		f(w, r)
	}
}

func ReplyJson(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		f(w, r)
	}
}

func ReplyMultipart(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "multipart/form-data")
		f(w, r)
	}
}

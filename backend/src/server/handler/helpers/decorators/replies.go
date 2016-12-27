package decorators

import "net/http"

func ReplyJson(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		f(w, r)
	}
}

func ReplyMultipart(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "multipart/form-data")
		f(w, r)
	}
}

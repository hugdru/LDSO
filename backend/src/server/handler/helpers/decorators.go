package helpers

import (
	"net/http"
	"server/handler/sessionData"
)

func RequestJson(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := GetContentType(r.Header.Get("Content-type"))
		if contentType != "application/json" {
			http.Error(w, Error("Expected, Content-type: application/json"), 415)
			return
		}
		f(w, r)
	}
}

func RequestMultipart(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := GetContentType(r.Header.Get("Content-type"))
		if contentType != "multipart/form-data" {
			http.Error(w, Error("Expected, Content-type: multipart/form-data"), 415)
			return
		}
		f(w, r)
	}
}

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

func OnlySuperadmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, "superadmin") {
			return
		}
		f(w, r)
	}
}

func OnlyLocaladmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, "localadmin") {
			return
		}
		f(w, r)
	}
}

func OnlyAuditors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, "auditor") {
			return
		}
		f(w, r)
	}
}

func OnlyClients(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, "client") {
			return
		}
		f(w, r)
	}
}

func OnlyWorkers(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRoles(r, []string{"superadmin", "localadmin", "auditor"}) {
			return
		}
		f(w, r)
	}
}

func permittedRole(r *http.Request, permittedRole string) bool {
	return permittedRoles(r, []string{permittedRole})
}

func permittedRoles(r *http.Request, permittedRoles []string) bool {
	entitySessionData, err := sessionData.GetSessionData(r)
	if err != nil {
		panic("Could not retrieve current SessionData")
	}
	for _, permittedRole := range permittedRoles {
		if entitySessionData.Role == permittedRole {
			return true
		}
	}
	return false
}
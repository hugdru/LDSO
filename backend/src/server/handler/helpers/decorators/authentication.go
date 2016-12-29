package decorators

import (
	"net/http"
	"server/handler/helpers"
	"server/handler/sessionData"
)

func OnlyLocaladminsOrAuditors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsLocaladminOrAuditor(r) {
			http.Error(w, helpers.Error("Only localadmins and auditors are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlySuperadmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsSuperadmin(r) {
			http.Error(w, helpers.Error("Only superadmins are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyLocaladmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsLocaladmin(r) {
			http.Error(w, helpers.Error("Only localadmins are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyAuditors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsAuditor(r) {
			http.Error(w, helpers.Error("Only auditors are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyClients(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsClient(r) {
			http.Error(w, helpers.Error("Only clients are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlySuperadminsOrLocaladminsOrClients(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsSuperadminOrLocaladminOrClient(r) {
			http.Error(w, helpers.Error("Only superadmins or localadmins or clients are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlySuperadminsOrLocaladmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsSuperadminOrLocaladmin(r) {
			http.Error(w, helpers.Error("Only superadmins or localadmins are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlySuperadminsOrLocaladminsOrAuditors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionData.IsSuperadminOrLocaladminOrAuditor(r) {
			http.Error(w, helpers.Error("Only superadmins or localadmins or auditors are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

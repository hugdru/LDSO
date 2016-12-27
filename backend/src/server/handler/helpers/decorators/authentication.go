package decorators

import (
	"net/http"
	"server/handler/sessionData"
	"server/handler/helpers"
)

func OnlySuperadmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, sessionData.Superadmin) {
			http.Error(w, helpers.Error("Only superadmins are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyLocaladmins(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, sessionData.Localadmin) {
			http.Error(w, helpers.Error("Only localadmins are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyAuditors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, sessionData.Auditor) {
			http.Error(w, helpers.Error("Only auditors are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyClients(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRole(r, sessionData.Client) {
			http.Error(w, helpers.Error("Only clients are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlyWorkers(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRoles(r, []string{sessionData.Superadmin, sessionData.Localadmin, sessionData.Auditor}) {
			http.Error(w, helpers.Error("Only superadmins, localadmins or auditors are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlySuperadminsOrClients(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRoles(r, []string{sessionData.Superadmin, sessionData.Client}) {
			http.Error(w, helpers.Error("Only superadmins or clients are permitted"), http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func OnlySuperadminsOrLocaladminsOrClients(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !permittedRoles(r, []string{sessionData.Superadmin, sessionData.Localadmin, sessionData.Client}) {
			http.Error(w, helpers.Error("Only superadmins or localadmins or clients are permitted"), http.StatusForbidden)
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

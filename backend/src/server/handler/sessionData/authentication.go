package sessionData

import "net/http"

func permittedRole(r *http.Request, permittedRole string) bool {
	return permittedRoles(r, []string{permittedRole})
}

func permittedRoles(r *http.Request, permittedRoles []string) bool {
	entitySessionData, err := GetSessionData(r)
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

func IsSuperadminOrLocaladmin(r *http.Request) bool {
	return permittedRoles(r, []string{Superadmin, Localadmin})
}

func IsSuperadminOrLocaladminOrAuditor(r *http.Request) bool {
	return permittedRoles(r, []string{Superadmin, Localadmin, Auditor})
}

func IsSuperadminOrLocaladminOrClient(r *http.Request) bool {
	return permittedRoles(r, []string{Superadmin, Localadmin, Client})
}

func IsClient(r *http.Request) bool {
	return permittedRole(r, Client)
}

func IsAuditor(r *http.Request) bool {
	return permittedRole(r, Auditor)
}

func IsLocaladmin(r *http.Request) bool {
	return permittedRole(r, Localadmin)
}

func IsSuperadmin(r *http.Request) bool {
	return permittedRole(r, Superadmin)
}

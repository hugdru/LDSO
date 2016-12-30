package sessionData

import (
	"encoding/gob"
	"github.com/alexedwards/scs/session"
	"net/http"
	"server/datastore"
)

const EntityKey = "entity"

const (
	Superadmin = "superadmin"
	Localadmin = "localadmin"
	Auditor    = "auditor"
	Client     = "client"
)

func GobRegister() {
	gob.Register(EntitySessionData{})
}

type EntitySessionData struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	IdCountry int64  `json:"idCountry"`
	Country   string `json:"country"`
	Role      string `json:"role"`
}

func GetSessionData(r *http.Request) (*EntitySessionData, error) {
	entitySessionData := &EntitySessionData{}
	err := session.GetObject(r, EntityKey, entitySessionData)
	if err != nil {
		return nil, err
	}
	return entitySessionData, err
}

func SetSessionData(entity *datastore.Entity, roleInterface interface{},
	w http.ResponseWriter, r *http.Request) *EntitySessionData {
	// Preventing session fixation
	err := session.RegenerateToken(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil
	}

	sessionData := &EntitySessionData{
		Id:        entity.Id,
		Username:  entity.Username,
		Email:     entity.Email,
		Name:      entity.Name,
		IdCountry: entity.IdCountry,
		Country:   entity.Country.Name,
		Role:      "",
	}

	//switch userRole := userRoleInterface.(type) {
	switch roleInterface.(type) {
	case *datastore.Superadmin:
		sessionData.Role = Superadmin
	case *datastore.Localadmin:
		sessionData.Role = Localadmin
	case *datastore.Auditor:
		sessionData.Role = Auditor
	case *datastore.Client:
		sessionData.Role = Client
	default:
		return nil
	}

	err = session.PutObject(r, EntityKey, sessionData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil
	}
	return sessionData
}

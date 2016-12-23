package sessionData

import (
	"encoding/gob"
	"github.com/alexedwards/scs/session"
	"net/http"
	"server/datastore"
)

const EntityKey = "entity"

func GobRegister() {
	gob.Register(SessionData{})
}

type SessionData struct {
	Id       int64
	Username string
	Email    string
	Name     string
	IdCountry int64
	Country  string
}

func SetSessionData(entity *datastore.Entity, w http.ResponseWriter, r *http.Request) *SessionData {
	// Preventing session fixation
	err := session.RegenerateToken(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil
	}

	sessionData := &SessionData{
		Id:       entity.Id,
		Username: entity.Username,
		Email:    entity.Email,
		Name:     entity.Username,
		IdCountry: entity.IdCountry,
		Country:  entity.Country.Name}

	err = session.PutObject(r, EntityKey, sessionData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil
	}
	return sessionData
}

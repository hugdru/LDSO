package handler

import (
	"github.com/pressly/chi"
	"server/handler/helpers"
	"net/http"
)

func (h *Handler) entiysRoutes(router chi.Router) {
	router.Get("/:id/login", helpers.ReplyJson(h.Login))
	router.Get("/:id/logout", helpers.ReplyJson(h.Logout))

}

//http://blog.brainattica.com/restful-json-api-jwt-go/
//https://gist.github.com/mschoebel/9398202

/*
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
*/
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		_,err := h.Datastore.GetEntityByUsernamePassword(name,pass)
		if err != nil{
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
		setSession(name, w)
		redirectTarget = "/internal"
	}
	http.Redirect(w, r, redirectTarget, 302)

}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	//uso de codificador
	/*if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}*/
	cookie := &http.Cookie{
		Name:  "session",
		Value: value["name"],
		Path:  "/",
	}
	http.SetCookie(response, cookie)

}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}


func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}


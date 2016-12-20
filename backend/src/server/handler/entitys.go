package handler

import (
	"github.com/pressly/chi"
	"server/handler/helpers"
	"net/http"
	"fmt"
	"server/datastore"

	"gopkg.in/guregu/null.v3/zero"
)

func (h *Handler) entiysRoutes(router chi.Router) {
	router.Get("/:id/login", helpers.ReplyJson(h.login))
	router.Get("/:id/logout", helpers.ReplyJson(h.logout))
	router.Get("/register", helpers.ReplyJson(h.register))

}

//http://blog.brainattica.com/restful-json-api-jwt-go/
//https://gist.github.com/mschoebel/9398202

/*
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
*/
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
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


func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	/*
	IdCountry   int64       `json:"idCountry" db:"id_country"`
	Name        string      `json:"name" db:"name"`
	Email       string      `json:"email" db:"email"`
	Username    string      `json:"username" db:"username"`
	Password    string      `json:"-" db:"password"`
	Image       []byte      `json:"image" db:"image"`
	Banned      zero.Bool   `json:"banned" db:"banned"`
	BannedDate  zero.Time   `json:"bannedDate" db:"banned_date"`
	Reason      zero.String `json:"reason" db:"reason"`
	Mobilephone zero.String `json:"mobilephone" db:"mobilephone"`
	Telephone   zero.String `json:"telephone" db:"telephone"`
	CreatedDate time.Time   `json:"createdDate" db:"created_date"`

	*/

	name := r.FormValue("name")
	email := r.FormValue("email")
	username := r.FormValue("username")
	//idCountry := r.FormValue("idCountry")
	//image := r.FormValue("image")
	mobilephone := r.FormValue("mobilephone")
	telephone := r.FormValue("telephone")
	pass := r.FormValue("password")
	userType := r.FormValue("type")
	//verificar o user name e email
	entity,err := h.Datastore.CheckEntityUsername(username,email)
	if name!="" && nil!=err && pass!="" {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	//entity.Country = idCountry
	entity.Mobilephone = zero.StringFrom(mobilephone)
	entity.Telephone =  zero.StringFrom(telephone)
	entity.Password = pass


	switch userType {
	case "auditor":
		fmt.Println("Create Audior ")
		auditor := datastore.NewAuditor(false)
		auditor.IdEntity = entity.Id
		err = h.Datastore.SaveAuditor(auditor)
		if nil!=err  {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	case "client":
		fmt.Println("Create Client")
		client := datastore.NewClient(false)
		client.IdEntity = entity.Id
		err = h.Datastore.SaveClient(client)
		if nil!=err  {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	case "localAdmin":
		fmt.Println("Create localAdmin")
		localAdmin := datastore.NewLocalAdmin(false)
		localAdmin.IdEntity = entity.Id
		err = h.Datastore.SaveLocalAdmin(localAdmin)
		if nil!=err  {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	default:
		fmt.Printf("That type does not exist")
		return
	}
	err = h.Datastore.SaveEntity(entity)
	if nil!=err  {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}



}

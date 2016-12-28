package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/alexedwards/scs/session"
	"github.com/elithrar/simple-scrypt"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"server/handler/sessionData"
	"strconv"
)

func (h *Handler) entitiesRoutes(router chi.Router) {
	router.Post("/login", decorators.ReplyJson(h.login))
	router.Get("/logout", decorators.ReplyJson(h.logout))
	router.Post("/register", decorators.ReplyJson(h.register))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		if decoder == nil {
			http.Error(w, helpers.Error("JSON decoder failed"), 500)
			return
		}
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
	case "multipart/form-data":
		input.Username = r.FormValue("username")
		input.Email = r.FormValue("email")
		input.Password = r.FormValue("password")
	default:
		http.Error(w, helpers.Error("JSON decoder failed"), 415)
		return
	}

	if (input.Username == "" && input.Email == "") || input.Password == "" {
		http.Error(w, helpers.Error("(Username or Email) and password"), 400)
		return
	}

	filter := make(map[string]interface{})
	if input.Username != "" {
		filter["username"] = input.Username
	}

	if input.Email != "" {
		filter["email"] = input.Email
	}

	if input.Password != "" {
	}

	entity, err := h.Datastore.GetEntityWithForeign(filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, helpers.Error("Bad credentials"), 401)
		return
	}

	var role interface{}

	superadmin, err := h.Datastore.GetSuperadminById(entity.Id)
	if err == nil {
		role = superadmin
	} else if err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		localadmin, err := h.Datastore.GetLocaladminById(entity.Id)
		if err == nil {
			role = localadmin
		} else if err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			auditor, err := h.Datastore.GetAuditorById(entity.Id)
			if err == nil {
				role = auditor
			} else if err != sql.ErrNoRows {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				client, err := h.Datastore.GetClientById(entity.Id)
				if err == nil {
					role = client
				} else if err != sql.ErrNoRows {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					http.Error(w, helpers.Error("This user has no role"), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	sessionData := sessionData.SetSessionData(entity, role, w, r)
	if sessionData == nil {
		http.Error(w, helpers.Error("Failed setting session"), 500)
		return
	}

	sessionDataSlice, err := json.Marshal(sessionData)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(sessionDataSlice)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	session.Destroy(w, r)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	contentLength := helpers.GetContentLength(r.Header.Get("Content-Length"))
	if contentLength == -1 {
		http.Error(w, helpers.Error("Invalid Content-Length header value"), http.StatusBadRequest)
		return
	}

	if contentLength > helpers.MaxMultipartSize {
		http.Error(w, helpers.Error("Data too big"), http.StatusBadRequest)
		return
	}

	var input struct {
		IdCountry     int64  `json:"idCountry"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		Username      string `json:"username"`
		Password      string `json:"password"`
		Mobilephone   string `json:"mobilephone"`
		Telephone     string `json:"telephone"`
		Role          string `json:"role"`
		imageBytes    []byte
		imageMimetype string
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		postIdCountry, err := helpers.ParseInt64(r.PostFormValue("idCountry"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.IdCountry = postIdCountry
		input.Name = r.PostFormValue("name")
		input.Email = r.PostFormValue("email")
		input.Username = r.PostFormValue("username")
		input.Password = r.PostFormValue("password")
		input.Mobilephone = r.PostFormValue("mobilephone")
		input.Telephone = r.PostFormValue("telephone")
		input.Role = r.PostFormValue("role")

		input.imageBytes, input.imageMimetype, err = helpers.ReadImage(r, "image", helpers.MaxImageFileSize)
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	default:
		http.Error(w, helpers.Error("Content-type not supported"), 415)
		return
	}

	if input.IdCountry == 0 || input.Name == "" || input.Email == "" || input.Username == "" || input.Password == "" || input.Role == "" {
		http.Error(w, helpers.Error("country, name, email, username, password, role are required"), 400)
		return
	}

	filter := make(map[string]interface{})
	if input.Username != "" {
		filter["username"] = input.Username
	}

	if input.Email != "" {
		filter["email"] = input.Email
	}

	err := h.Datastore.CheckEntityExists(filter)
	if err == nil {
		http.Error(w, helpers.Error("already Exists"), 400)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	entity := datastore.NewEntity(false)
	entity.IdCountry = input.IdCountry
	entity.Name = input.Name
	entity.Email = input.Email
	entity.Username = input.Username
	hash, err := scrypt.GenerateFromPassword([]byte(input.Password), scrypt.DefaultParams)
	if err != nil {
		http.Error(w, helpers.Error("Failed to create hash"), 500)
		return
	}
	entity.Password = string(hash)
	entity.Mobilephone = zero.StringFrom(input.Mobilephone)
	entity.Telephone = zero.StringFrom(input.Telephone)
	entity.ImageMimetype = zero.StringFrom(input.imageMimetype)
	entity.Image = input.imageBytes
	entity.CreatedDate = helpers.TheTime()

	err = h.Datastore.SaveEntity(entity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	switch input.Role {
	case sessionData.Client:
		client := datastore.NewClient(false)
		client.IdEntity = entity.Id
		err = h.Datastore.SaveClient(client)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	case sessionData.Auditor:
		if !sessionData.IsSuperadminOrLocaladmin(r) {
			http.Error(w, helpers.Error("Not superadmin or localadmin"), http.StatusForbidden)
			return
		}
		auditor := datastore.NewAuditor(false)
		auditor.IdEntity = entity.Id
		err = h.Datastore.SaveAuditor(auditor)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	case sessionData.Localadmin:
		if !sessionData.IsSuperadmin(r) {
			http.Error(w, helpers.Error("Not superadmin"), http.StatusForbidden)
			return
		}
		localadmin := datastore.NewLocaladmin(false)
		localadmin.IdEntity = entity.Id
		err = h.Datastore.SaveLocaladmin(localadmin)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	case sessionData.Superadmin:
		if !sessionData.IsSuperadmin(r) {
			http.Error(w, helpers.Error("Not superadmin"), http.StatusForbidden)
			return
		}
		superadmin := datastore.NewSuperadmin(false)
		superadmin.IdEntity = entity.Id
		err = h.Datastore.SaveSuperadmin(superadmin)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	default:
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	response := []byte(`{"id":"` + strconv.FormatInt(entity.Id, 10) + "}")
	w.Write(response)
}

package handler

import (
	"context"
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

func (h *Handler) entityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idEntity, err := helpers.ParseInt64(chi.URLParam(r, "ide"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		entity, err := h.Datastore.GetClientById(idEntity, false, true)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "client", entity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

	entity, err := h.Datastore.GetEntity(filter, true, false)
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

	superadmin, err := h.Datastore.GetSuperadminById(entity.Id, false, false)
	if err == nil {
		role = superadmin
	} else if err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		localadmin, err := h.Datastore.GetLocaladminById(entity.Id, false, false)
		if err == nil {
			role = localadmin
		} else if err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			auditor, err := h.Datastore.GetAuditorById(entity.Id, false, false)
			if err == nil {
				role = auditor
			} else if err != sql.ErrNoRows {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				client, err := h.Datastore.GetClientById(entity.Id, false, false)
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
		http.Error(w, helpers.Error("Failed setting session"), http.StatusInternalServerError)
		return
	}

	sessionDataSlice, err := json.Marshal(sessionData)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(sessionDataSlice)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	err := session.Destroy(w, r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	tx, successful := func() (*sql.Tx, bool) {
		contentLength := helpers.GetContentLength(r.Header.Get("Content-Length"))
		if contentLength == -1 {
			http.Error(w, helpers.Error("Invalid Content-Length header value"), http.StatusBadRequest)
			return nil, false
		}

		if contentLength > helpers.MaxMultipartSize {
			http.Error(w, helpers.Error("Data too big"), http.StatusBadRequest)
			return nil, false
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
			Image         string `json:"image"`
			imageBytes    []byte `json:"-"`
			imageMimetype string `json:"-"`
			imageHash     string `json:"-"`
		}

		contentType := helpers.GetContentType(r.Header.Get("Content-type"))
		switch contentType {
		case "application/json":
			decoder := json.NewDecoder(r.Body)
			if decoder == nil {
				http.Error(w, helpers.Error("JSON decoder failed"), http.StatusInternalServerError)
				return nil, false
			}
			err := decoder.Decode(&input)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), http.StatusBadRequest)
				return nil, false
			}
			if input.Image != "" {
				input.imageBytes, input.imageMimetype, input.imageHash, err = helpers.ReadImageBase64(input.Image, helpers.MaxImageFileSize)
				if err != nil {
					http.Error(w, helpers.Error(err.Error()), http.StatusInternalServerError)
					return nil, false
				}
			}
		case "multipart/form-data":
			var err error
			input.IdCountry, _ = helpers.ParseInt64(r.PostFormValue("idCountry"))
			input.Name = r.PostFormValue("name")
			input.Email = r.PostFormValue("email")
			input.Username = r.PostFormValue("username")
			input.Password = r.PostFormValue("password")
			input.Mobilephone = r.PostFormValue("mobilephone")
			input.Telephone = r.PostFormValue("telephone")
			input.Role = r.PostFormValue("role")

			input.imageBytes, input.imageMimetype, input.imageHash, err = helpers.ReadImage(r, "image", helpers.MaxImageFileSize)
			if err != nil && err != http.ErrMissingFile {
				http.Error(w, helpers.Error(err.Error()), 500)
				return nil, false
			}
		default:
			http.Error(w, helpers.Error("Content-type not supported"), 415)
			return nil, false
		}

		if input.Name == "" || input.Email == "" || input.Username == "" || input.Password == "" || input.Role == "" {
			http.Error(w, helpers.Error("country, name, email, username, password, role are required"), 400)
			return nil, false
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
			return nil, false
		} else if err != sql.ErrNoRows {
			http.Error(w, helpers.Error(err.Error()), 400)
			return nil, false
		}

		entity := datastore.NewEntity(false)
		if input.IdCountry == 0 {
			portugal, err := h.Datastore.GetCountryByName("Portugal")
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return nil, false
			}
			entity.IdCountry = portugal.Id
		} else {
			entity.IdCountry = input.IdCountry
		}
		entity.Name = input.Name
		entity.Email = input.Email
		entity.Username = input.Username
		hash, err := scrypt.GenerateFromPassword([]byte(input.Password), scrypt.DefaultParams)
		if err != nil {
			http.Error(w, helpers.Error("Failed to create hash"), 500)
			return nil, false
		}
		entity.Password = string(hash)
		entity.Mobilephone = zero.StringFrom(input.Mobilephone)
		entity.Telephone = zero.StringFrom(input.Telephone)
		entity.ImageMimetype = zero.StringFrom(input.imageMimetype)
		entity.Image = input.imageBytes
		entity.ImageHash = zero.StringFrom(input.imageHash)
		entity.CreatedDate = helpers.TheTime()

		tx, err := h.Datastore.BeginTransaction()
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return nil, false
		}

		err = h.Datastore.SaveEntityTx(tx, entity)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return tx, false
		}

		switch input.Role {
		case sessionData.Client:
			client := datastore.NewClient(false)
			client.IdEntity = entity.Id
			err = h.Datastore.SaveClientTx(tx, client)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return tx, false
			}
		case sessionData.Auditor:
			if !sessionData.IsSuperadminOrLocaladmin(r) {
				http.Error(w, helpers.Error("Not superadmin or localadmin"), http.StatusForbidden)
				return tx, false
			}
			auditor := datastore.NewAuditor(false)
			auditor.IdEntity = entity.Id
			err = h.Datastore.SaveAuditorTx(tx, auditor)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return tx, false
			}
		case sessionData.Localadmin:
			if !sessionData.IsSuperadmin(r) {
				http.Error(w, helpers.Error("Not superadmin"), http.StatusForbidden)
				return tx, false
			}
			localadmin := datastore.NewLocaladmin(false)
			localadmin.IdEntity = entity.Id
			err = h.Datastore.SaveLocaladminTx(tx, localadmin)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return tx, false
			}
		case sessionData.Superadmin:
			if !sessionData.IsSuperadmin(r) {
				http.Error(w, helpers.Error("Not superadmin"), http.StatusForbidden)
				return tx, false
			}
			superadmin := datastore.NewSuperadmin(false)
			superadmin.IdEntity = entity.Id
			err = h.Datastore.SaveSuperadminTx(tx, superadmin)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return tx, false
			}
		default:
			http.Error(w, helpers.Error(err.Error()), 500)
			return tx, false
		}

		response := []byte(`{"id":` + strconv.FormatInt(entity.Id, 10) + `}`)
		w.Write(response)

		return tx, true
	}()
	if tx != nil {
		if successful {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}

func getEntityImage(w http.ResponseWriter, r *http.Request, entity *datastore.Entity, urlHash, redirectUrl string) bool {
	if entity == nil {
		http.Error(w, helpers.Error("entity != nil"), http.StatusInternalServerError)
		return false
	}

	if entity.ImageMimetype.Valid && entity.Image != nil && entity.ImageHash.Valid {
		if err := helpers.ImageCashingControlWriter(w, r, urlHash, entity.ImageHash.String, entity.ImageMimetype.String, entity.Image, redirectUrl); err != nil {
			http.Error(w, helpers.Error(err.Error()), http.StatusInternalServerError)
			return false
		}
		return true
	}
	http.Error(w, helpers.Error("no image"), http.StatusNotFound)
	return false
}

func (h *Handler) updateEntity(w http.ResponseWriter, r *http.Request, tx *sql.Tx, entity *datastore.Entity) bool {
	if entity == nil {
		http.Error(w, helpers.Error("entity != nil"), http.StatusInternalServerError)
		return false
	}

	var input struct {
		Entity struct {
			IdCountry     int64  `json:"idCountry"`
			Name          string `json:"name"`
			Email         string `json:"email"`
			Username      string `json:"username"`
			Password      string `json:"password"`
			Mobilephone   string `json:"mobilephone"`
			Telephone     string `json:"telephone"`
			Image         string `json:"image"`
			imageBytes    []byte `json:"-"`
			imageMimetype string `json:"-"`
			imageHash     string `json:"-"`
		} `json:"entity"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		if decoder == nil {
			http.Error(w, helpers.Error("JSON decoder failed"), http.StatusInternalServerError)
			return false
		}
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), http.StatusBadRequest)
			return false
		}
		if input.Entity.Image != "" {
			input.Entity.imageBytes, input.Entity.imageMimetype, input.Entity.imageHash, err = helpers.ReadImageBase64(input.Entity.Image, helpers.MaxImageFileSize)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), http.StatusInternalServerError)
				return false
			}
		}
	case "multipart/form-data":
		var err error
		input.Entity.IdCountry, _ = helpers.ParseInt64(r.PostFormValue("entity.idCountry"))
		input.Entity.Name = r.PostFormValue("entity.name")
		input.Entity.Email = r.PostFormValue("entity.email")
		input.Entity.Username = r.PostFormValue("entity.username")
		input.Entity.Password = r.PostFormValue("entity.password")
		input.Entity.Mobilephone = r.PostFormValue("entity.mobilephone")
		input.Entity.Telephone = r.PostFormValue("entity.telephone")

		input.Entity.imageBytes, input.Entity.imageMimetype, input.Entity.imageHash, err = helpers.ReadImage(r, "image", helpers.MaxImageFileSize)
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, helpers.Error(err.Error()), 500)
			return false
		}
	default:
		http.Error(w, helpers.Error("Content-type not supported"), 415)
		return false
	}

	if input.Entity.IdCountry == 0 && input.Entity.Name == "" && input.Entity.Email == "" && input.Entity.Username == "" && input.Entity.Password == "" && input.Entity.Mobilephone == "" && input.Entity.Telephone == "" && input.Entity.imageBytes == nil {
		http.Error(w, helpers.Error("At least one"), 400)
		return false
	}

	if input.Entity.IdCountry != 0 {
		entity.IdCountry = input.Entity.IdCountry
	}

	if input.Entity.Name != "" {
		entity.Name = input.Entity.Name
	}

	if input.Entity.Email != "" {
		entity.Email = input.Entity.Email
	}

	if input.Entity.Username != "" {
		entity.Username = input.Entity.Username
	}

	if input.Entity.Password != "" {
		hash, err := scrypt.GenerateFromPassword([]byte(input.Entity.Password), scrypt.DefaultParams)
		if err != nil {
			http.Error(w, helpers.Error("Failed to create hash"), 500)
			return false
		}
		entity.Password = string(hash)
	}

	if input.Entity.imageBytes != nil {
		entity.Image = input.Entity.imageBytes
		entity.ImageMimetype = zero.StringFrom(input.Entity.imageMimetype)
		entity.ImageHash = zero.StringFrom(input.Entity.imageHash)
	}

	if input.Entity.Mobilephone != "" {
		entity.Mobilephone = zero.StringFrom(input.Entity.Mobilephone)
	}

	if input.Entity.Telephone != "" {
		entity.Telephone = zero.StringFrom(input.Entity.Telephone)
	}

	err := h.Datastore.SaveEntityTx(tx, entity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return false
	}
	return true
}

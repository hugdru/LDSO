package handler

import (
	"github.com/pressly/chi"
	"server/handler/helpers"
	"net/http"
	"server/datastore"
	"gopkg.in/guregu/null.v3/zero"
	"encoding/json"
	"github.com/elithrar/simple-scrypt"
	"encoding/gob"
	"github.com/alexedwards/scs/session"
	"strconv"
	"io/ioutil"
	"time"
)

const maxMultipartSize = 32 << 20
const maxImageFileSize = 16 << 20

type CookieData struct {
	id  int64
	username string
	email string
	name string
	country string
}

var acceptedImageTypes = map[string] bool {
	"image/gif": true,
	"image/png": true,
	"image/jpeg": true,
}

const EntityKey = "entity"

func (h *Handler) entitiesRoutes(router chi.Router) {
	gob.Register(CookieData{})
	router.Get("/login", helpers.ReplyJson(h.login))
	router.Get("/logout", helpers.ReplyJson(h.logout))
	router.Get("/register", helpers.ReplyJson(h.register))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Username string `json:"username"`
		Email  string `json:"email"`
		Password string `json:"password"`
	}

	switch(r.Header.Get("Content-type")) {
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
		input.Username = r.FormValue("username");
		input.Email = r.FormValue("email");
		input.Password = r.FormValue("password");
	default:
		http.Error(w, helpers.Error("JSON decoder failed"), 415)
		return
	}

	if (input.Username == "" && input.Email == "") || input.Password == "" {
		http.Error(w, helpers.Error("(Username or Email) and password"), 400)
		return
	}

	filter := make(map[string]string)
	if input.Username != "" {
		filter["username"] = input.Username;
	}

	if input.Email != "" {
		filter["email"] = input.Email;
	}

	if input.Password != "" {
	}

	entity, err := h.Datastore.GetEntity(filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, helpers.Error("Bad credentials"), 401)
		return
	}

	if !setSession(entity, w, r) {
		return;
	}
}

func setSession(entity *datastore.Entity, w http.ResponseWriter, r *http.Request) bool {
	// Preventing session fixation
	err := session.RegenerateToken(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return false;
	}

	cookieData := &CookieData{
		id: entity.Id,
		username: entity.Username,
		email: entity.Email,
		name: entity.Username,
		country: entity.Country.Name}

	err = session.PutObject(r, EntityKey, cookieData);
	if err != nil {
		http.Error(w, err.Error(), 500)
		return false;
	}
	return true;
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	session.Destroy(w, r);
}
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	contentLengthStr := r.Header.Get("Content-Length")
	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if contentLength > maxMultipartSize {
		http.Error(w, helpers.Error("data too big"), 400)
		return
	}

	var input struct {
		IdCountry int64 `json:"idCountry"`
		Name string `json:"name"`
		Email string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Mobilephone string `json:"mobilephone"`
		Telephone   string `json:"telephone"`
		imageBytes []byte
		imageMime string
	}

	switch(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		postIdCountry, err := strconv.ParseInt(r.PostFormValue("idCountry"), 10, 64)
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
		mFile, mFileHeader, err := r.FormFile("image")
		if err == nil {
			defer mFile.Close()

			imageContentType := mFileHeader.Header.Get("Content-type")
			if !acceptedImageTypes[imageContentType] {
				http.Error(w, helpers.Error("Gif, jpeg or png image"), 400)
				return
			}
			mimetypeBuffer := make([]byte, 512)
			if _, err = mFile.Read(mimetypeBuffer); err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return
			}
			if (imageContentType != http.DetectContentType(mimetypeBuffer)) {
				http.Error(w, helpers.Error("Gif, jpeg or png image"), 400)
				return
			}
			input.imageMime = imageContentType;
			imageContentLengthStr := mFileHeader.Header.Get("Content-Length")
			imageContentLength, err := strconv.ParseInt(imageContentLengthStr, 10, 64)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 400)
				return
			}
			if imageContentLength > maxImageFileSize {
				http.Error(w, helpers.Error("image is too big, max: " + strconv.FormatInt(maxImageFileSize, 10)) + " bytes", 400)
				return
			}

			input.imageBytes, err = ioutil.ReadAll(mFile)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 400)
				return
			}
		} else if err != http.ErrMissingFile {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	default:
		http.Error(w, helpers.Error("Content-type not supported"), 415)
		return
	}

	if input.IdCountry == 0 || input.Name == "" || input.Email == "" || input.Username == "" || input.Password == "" {
		http.Error(w, helpers.Error("country, name, email, username, password are required"), 400)
		return
	}

	filter := make(map[string]string)
	if input.Username != "" {
		filter["username"] = input.Username;
	}

	if input.Email != "" {
		filter["email"] = input.Email;
	}

	exists, err := h.Datastore.CheckEntityExists(filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	if exists {
		http.Error(w, helpers.Error("Already exists"), 400)
		return
	}

	entity := datastore.NewEntity(false);
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
	entity.ImageMimetype = zero.StringFrom(input.imageMime)
	entity.Image = input.imageBytes
	entity.CreatedDate = time.Now().UTC()

	err = h.Datastore.SaveEntity(entity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

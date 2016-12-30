package handler

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"server/handler/sessionData"
)

func (h *Handler) localadminsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getLocaladmins))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createLocaladmin)))
	router.Route("/:idl", h.localadminRoutes)
}

func (h *Handler) localadminRoutes(router chi.Router) {
	router.Use(h.localadminContext)
	router.Get("/", decorators.ReplyJson(h.getLocaladmin))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateLocaladmin)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteLocaladmin)))
}

func (h *Handler) localadminContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idLocaladminStr := chi.URLParam(r, "idl")
		idLocaladmin, err := helpers.ParseInt64(idLocaladminStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		esd, err := sessionData.GetSessionData(r)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		restricted := true
		switch esd.Role {
		case sessionData.Superadmin:
			restricted = false
		case sessionData.Localadmin:
			if esd.Id == idLocaladmin {
				restricted = false
			}
		}

		localadmin, err := h.Datastore.GetLocaladminById(idLocaladmin, true, restricted)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		if r.Method != http.MethodGet {
			if esd.Id != localadmin.IdEntity && esd.Role != sessionData.Superadmin {
				http.Error(w, helpers.Error("Not the owner of the account"), http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), "localadmin", localadmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getLocaladmins(w http.ResponseWriter, r *http.Request) {

	esd, err := sessionData.GetSessionData(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	restricted := true
	switch esd.Role {
	case sessionData.Superadmin:
		restricted = false
	}

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	localadmins, err := h.Datastore.GetLocaladmins(limit, offset, true, restricted)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	localadminsSlice, err := json.Marshal(localadmins)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(localadminsSlice)
}

func (h *Handler) createLocaladmin(w http.ResponseWriter, r *http.Request) {

	var input struct {
		IdEntity int64 `json:"IdEntity"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdEntity, err = helpers.ParseInt64(r.PostFormValue("IdEntity"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
	case "application/json":
		d := json.NewDecoder(r.Body)
		err := d.Decode(&input)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
	default:
		http.Error(w, helpers.Error("Content-type not supported"), 415)
		return
	}

	newLocaladmin := datastore.NewLocaladmin(false)
	err := newLocaladmin.MustSet(input.IdEntity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveLocaladmin(newLocaladmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	localadminSlice, err := json.Marshal(newLocaladmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(localadminSlice)
}

func (h *Handler) getLocaladmin(w http.ResponseWriter, r *http.Request) {

	localadmin := r.Context().Value("localadmin").(*datastore.Localadmin)

	localadminSlice, err := json.Marshal(localadmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(localadminSlice)
}

func (h *Handler) updateLocaladmin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, helpers.Error("Not implemented because localadmin only has one column idEntity"), http.StatusNotImplemented)
}

func (h *Handler) deleteLocaladmin(w http.ResponseWriter, r *http.Request) {

	localadmin := r.Context().Value("localadmin").(*datastore.Localadmin)

	err := h.Datastore.DeleteLocaladmin(localadmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

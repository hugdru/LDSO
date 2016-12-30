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

func (h *Handler) superadminsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getSuperadmins))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createSuperadmin)))
	router.Route("/:ids", h.superadminRoutes)
}

func (h *Handler) superadminRoutes(router chi.Router) {
	router.Use(h.superadminContext)
	router.Get("/", decorators.ReplyJson(h.getSuperadmin))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateSuperadmin)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteSuperadmin)))
}

func (h *Handler) superadminContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idSuperadminStr := chi.URLParam(r, "idl")
		idSuperadmin, err := helpers.ParseInt64(idSuperadminStr)
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
		}

		superadmin, err := h.Datastore.GetSuperadminById(idSuperadmin, true, restricted)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		if r.Method != http.MethodGet {
			if esd.Role != sessionData.Superadmin {
				http.Error(w, helpers.Error("Not a superadmin"), http.StatusForbidden)
				return
			}
			if r.Method != http.MethodDelete {
				if esd.Id != idSuperadmin {
					http.Error(w, helpers.Error("Not the owner"), http.StatusForbidden)
					return
				}
			}
		}
		ctx := context.WithValue(r.Context(), "superadmin", superadmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getSuperadmins(w http.ResponseWriter, r *http.Request) {

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

	superadmins, err := h.Datastore.GetSuperadmins(limit, offset, true, restricted)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	superadminsSlice, err := json.Marshal(superadmins)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(superadminsSlice)
}

func (h *Handler) createSuperadmin(w http.ResponseWriter, r *http.Request) {

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

	newSuperadmin := datastore.NewSuperadmin(false)
	err := newSuperadmin.MustSet(input.IdEntity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveSuperadmin(newSuperadmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	superadminSlice, err := json.Marshal(newSuperadmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(superadminSlice)
}

func (h *Handler) getSuperadmin(w http.ResponseWriter, r *http.Request) {

	superadmin := r.Context().Value("superadmin").(*datastore.Superadmin)

	superadminSlice, err := json.Marshal(superadmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(superadminSlice)
}

func (h *Handler) updateSuperadmin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, helpers.Error("Not implemented because superadmin only has one column idEntity"), http.StatusNotImplemented)
}

func (h *Handler) deleteSuperadmin(w http.ResponseWriter, r *http.Request) {

	superadmin := r.Context().Value("superadmin").(*datastore.Superadmin)

	err := h.Datastore.DeleteSuperadmin(superadmin)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

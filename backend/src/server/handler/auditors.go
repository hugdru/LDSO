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

func (h *Handler) auditorsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getAuditors))
	router.Post("/", decorators.OnlySuperadminsOrLocaladmins(decorators.ReplyJson(h.createAuditor)))
	router.Route("/:ida", h.auditorRoutes)
}

func (h *Handler) auditorRoutes(router chi.Router) {
	router.Use(h.auditorContext)
	router.Get("/", decorators.ReplyJson(h.getAuditor))
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.updateAuditor)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteAuditor)))
}

func (h *Handler) auditorContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAuditorStr := chi.URLParam(r, "ida")
		idAuditor, err := helpers.ParseInt64(idAuditorStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		auditor, err := h.Datastore.GetAuditorByIdWithEntity(idAuditor)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		if r.Method != http.MethodGet {
			entitySessionData, err := sessionData.GetSessionData(r)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 400)
				return
			}
			if entitySessionData.Id != auditor.IdEntity && entitySessionData.Role != sessionData.Superadmin && entitySessionData.Role != sessionData.Localadmin {
				http.Error(w, helpers.Error("Not the owner of the account"), http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), "auditor", auditor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getAuditors(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditors, err := h.Datastore.GetAuditorsWithEntity(limit, offset)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditorsSlice, err := json.Marshal(auditors)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditorsSlice)
}

func (h *Handler) createAuditor(w http.ResponseWriter, r *http.Request) {

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

	newAuditor := datastore.NewAuditor(false)
	err := newAuditor.MustSet(input.IdEntity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveAuditor(newAuditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	auditorSlice, err := json.Marshal(newAuditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(auditorSlice)
}

func (h *Handler) getAuditor(w http.ResponseWriter, r *http.Request) {

	auditor := r.Context().Value("auditor").(*datastore.Auditor)

	auditorSlice, err := json.Marshal(auditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditorSlice)
}

func (h *Handler) updateAuditor(w http.ResponseWriter, r *http.Request) {
	http.Error(w, helpers.Error("Not implemented because auditor only has one column idEntity"), http.StatusNotImplemented)
}

func (h *Handler) deleteAuditor(w http.ResponseWriter, r *http.Request) {

	auditor := r.Context().Value("auditor").(*datastore.Auditor)

	err := h.Datastore.DeleteAuditor(auditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

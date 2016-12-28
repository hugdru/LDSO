package handler

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
)

func (h *Handler) maingroupsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getMaingroups))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createMaingroup)))
	router.Route("/:idm", h.maingroupRoutes)
}

func (h *Handler) maingroupRoutes(router chi.Router) {
	router.Use(h.maingroupContext)
	router.Get("/", decorators.ReplyJson(h.getMaingroup))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateMaingroup)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteMaingroup)))
}

func (h *Handler) maingroupContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idMaingroupStr := chi.URLParam(r, "idm")
		idMaingroup, err := helpers.ParseInt64(idMaingroupStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		maingroup, err := h.Datastore.GetMaingroupById(idMaingroup)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "maingroup", maingroup)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getMaingroups(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"idTemplate", "id_template"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	maingroups, err := h.Datastore.GetMaingroups(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	maingroupsSlice, err := json.Marshal(maingroups)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(maingroupsSlice)
}

func (h *Handler) createMaingroup(w http.ResponseWriter, r *http.Request) {

	var input struct {
		IdTemplate int64  `json:"idTemplate"`
		Name       string `json:"name"`
		Weight     int    `json:"weight"`
	}
	input.Weight = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdTemplate, err = helpers.ParseInt64(r.PostFormValue("idTemplate"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.Name = r.PostFormValue("name")
		input.Weight, err = helpers.ParseInt(r.PostFormValue("weight"))
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

	maingroup := datastore.NewMaingroup(false)
	err := maingroup.MustSet(input.IdTemplate, input.Name, input.Weight)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	maingroup.CreatedDate = helpers.TheTime()

	err = h.Datastore.SaveMaingroup(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	maingroupSlice, err := json.Marshal(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(maingroupSlice)
}

func (h *Handler) getMaingroup(w http.ResponseWriter, r *http.Request) {

	maingroup := r.Context().Value("maingroup").(*datastore.Maingroup)

	maingroupSlice, err := json.Marshal(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(maingroupSlice)
}

func (h *Handler) updateMaingroup(w http.ResponseWriter, r *http.Request) {

	maingroup := r.Context().Value("maingroup").(*datastore.Maingroup)

	var input struct {
		Name   string `json:"name"`
		Weight int    `json:"weight"`
	}
	input.Weight = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.Name = r.PostFormValue("name")
		input.Weight, err = helpers.ParseInt(r.PostFormValue("weight"))
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

	if input.Name == "" && input.Weight == -1 {
		http.Error(w, helpers.Error("[name] [weight]"), 400)
		return
	}

	err := maingroup.UpdateSetIfNotEmptyOrNil(input.Name, input.Weight)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveMaingroup(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteMaingroup(w http.ResponseWriter, r *http.Request) {

	maingroup := r.Context().Value("maingroup").(*datastore.Maingroup)

	err := h.Datastore.DeleteMaingroup(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

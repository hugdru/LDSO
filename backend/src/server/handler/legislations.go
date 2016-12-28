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

func (h *Handler) legislationsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getLegislations))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createLegislation)))
	router.Route("/:idl", h.legislationRoutes)
}

func (h *Handler) legislationRoutes(router chi.Router) {
	router.Use(h.legislationContext)
	router.Get("/", decorators.ReplyJson(h.getLegislation))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateLegislation)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteLegislation)))
}

func (h *Handler) legislationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idLegislationStr := chi.URLParam(r, "idl")
		idLegislation, err := helpers.ParseInt64(idLegislationStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		legislation, err := h.Datastore.GetLegislationById(idLegislation)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "legislation", legislation)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getLegislations(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"name"},
		[]string{"description"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	legislations, err := h.Datastore.GetLegislations(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	legislationsSlice, err := json.Marshal(legislations)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(legislationsSlice)
}

func (h *Handler) createLegislation(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Url         string `json:"url"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.Name = r.PostFormValue("name")
		input.Description = r.PostFormValue("description")
		input.Url = r.PostFormValue("url")
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

	legislation := datastore.NewLegislation(false)
	err := legislation.MustSet(input.Name)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = legislation.OptionalSetIfNotEmptyOrNil(input.Description, input.Url)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveLegislation(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	legislationSlice, err := json.Marshal(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(legislationSlice)
}

func (h *Handler) getLegislation(w http.ResponseWriter, r *http.Request) {

	legislation := r.Context().Value("legislation").(*datastore.Legislation)

	legislationSlice, err := json.Marshal(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(legislationSlice)
}

func (h *Handler) updateLegislation(w http.ResponseWriter, r *http.Request) {

	legislation := r.Context().Value("legislation").(*datastore.Legislation)

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Url         string `json:"url"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.Name = r.PostFormValue("name")
		input.Description = r.PostFormValue("description")
		input.Url = r.PostFormValue("url")
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

	if input.Name == "" && input.Description == "" && input.Url == "" {
		http.Error(w, helpers.Error("At least one of: [name] [description] [url]"), 400)
		return
	}

	err := legislation.AllSetIfNotEmptyOrNil(input.Name, input.Description, input.Url)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveLegislation(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteLegislation(w http.ResponseWriter, r *http.Request) {

	legislation := r.Context().Value("legislation").(*datastore.Legislation)

	err := h.Datastore.DeleteLegislation(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

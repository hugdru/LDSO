package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) legislationsRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getLegislations))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createLegislation)))
	router.Get("/:id", helpers.ReplyJson(h.getLegislation))
	router.Put("/:id", helpers.RequestJson(helpers.ReplyJson(h.updateLegislation)))
	router.Delete("/:id", helpers.ReplyJson(h.deleteLegislation))
}

func (h *Handler) getLegislations(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"}, []string{"name"}, []string{"description"},
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

func (h *Handler) getLegislation(w http.ResponseWriter, r *http.Request) {
	idLegislation := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idLegislation, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	legislation, err := h.Datastore.GetLegislationById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	legislationSlice, err := json.Marshal(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(legislationSlice)
}

func (h *Handler) createLegislation(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Name        string
		Description string
		Url         string
	}

	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" {
		http.Error(w, helpers.Error("name [description] [url]"), 400)
		return
	}

	legislation := datastore.NewLegislation(true)
	legislation.Name = input.Name
	legislation.Description = zero.StringFrom(input.Description)
	legislation.Url = zero.StringFrom(input.Url)

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

func (h *Handler) updateLegislation(w http.ResponseWriter, r *http.Request) {
	idLegislation := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idLegislation, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	legislation, err := h.Datastore.GetLegislationById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	d := json.NewDecoder(r.Body)
	if d == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Name        string
		Description string
		Url         string
	}

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" && input.Description == "" && input.Url == "" {
		http.Error(w, helpers.Error("At least one of: [name] [description] [url]"), 400)
		return
	}

	if input.Name != "" {
		legislation.Name = input.Name
	}

	if input.Description != "" {
		legislation.Description = zero.StringFrom(input.Description)
	}

	if input.Url != "" {
		legislation.Url = zero.StringFrom(input.Url)
	}

	err = h.Datastore.SaveLegislation(legislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteLegislation(w http.ResponseWriter, r *http.Request) {
	idLegislation := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idLegislation, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteLegislationById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

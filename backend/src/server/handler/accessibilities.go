package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) accessibilitiesRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getAccessibilities))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createAccessibility)))
	router.Get("/:id", helpers.ReplyJson(h.getAccessibility))
	router.Put("/:id", helpers.RequestJson(helpers.ReplyJson(h.updateAccessibility)))
	router.Delete("/:id", helpers.ReplyJson(h.deleteAccessibility))
}

func (h *Handler) getAccessibilities(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	accessibilities, err := h.Datastore.GetAccessibilities(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	accessibilitiesSlice, err := json.Marshal(accessibilities)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(accessibilitiesSlice)
}

func (h *Handler) getAccessibility(w http.ResponseWriter, r *http.Request) {
	idAccessibility := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAccessibility, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	accessibility, err := h.Datastore.GetAccessibilityById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	accessibilitySlice, err := json.Marshal(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(accessibilitySlice)
}

func (h *Handler) createAccessibility(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Name string
	}

	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" {
		http.Error(w, helpers.Error("name"), 400)
		return
	}

	accessibility := datastore.NewAccessibility(true)
	accessibility.Name = input.Name

	err = h.Datastore.SaveAccessibility(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	accessibilitySlice, err := json.Marshal(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(accessibilitySlice)
}

func (h *Handler) updateAccessibility(w http.ResponseWriter, r *http.Request) {
	idAccessibility := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAccessibility, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	accessibility, err := h.Datastore.GetAccessibilityById(id)
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
		Name string
	}

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" {
		http.Error(w, helpers.Error("name"), 400)
		return
	}

	if input.Name != "" {
		accessibility.Name = input.Name
	}

	err = h.Datastore.SaveAccessibility(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteAccessibility(w http.ResponseWriter, r *http.Request) {
	idAccessibility := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAccessibility, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAccessibilityById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

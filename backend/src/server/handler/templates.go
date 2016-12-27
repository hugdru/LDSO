package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"
	"server/handler/helpers/decorators"
)

func (h *Handler) templatesRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getTemplates))
	router.Post("/", decorators.RequestJson(decorators.ReplyJson(h.createTemplate)))
	router.Get("/:id", decorators.ReplyJson(h.getTemplate))
	router.Put("/:id", decorators.RequestJson(decorators.ReplyJson(h.updateTemplate)))
	router.Delete("/:id", decorators.ReplyJson(h.deleteTemplate))
}

func (h *Handler) getTemplates(w http.ResponseWriter, r *http.Request) {

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

	templates, err := h.Datastore.GetTemplates(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	templatesSlice, err := json.Marshal(templates)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(templatesSlice)
}

func (h *Handler) getTemplate(w http.ResponseWriter, r *http.Request) {
	idTemplate := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idTemplate, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	template, err := h.Datastore.GetTemplate(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	templateSlice, err := json.Marshal(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(templateSlice)
}

func (h *Handler) createTemplate(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	if d == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Name        string
		Description string
	}

	err := d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" {
		http.Error(w, helpers.Error("The template must have a name"), 400)
		return
	}

	template := datastore.NewTemplate(false)
	template.Name = input.Name
	template.Description = zero.StringFrom(input.Description)
	template.CreatedDate = helpers.TheTime()

	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = h.Datastore.SaveTemplate(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	templateSlice, err := json.Marshal(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(templateSlice)
}

func (h *Handler) updateTemplate(w http.ResponseWriter, r *http.Request) {
	idTemplate := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idTemplate, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	template, err := h.Datastore.GetTemplateById(id)
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
	}

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" && input.Description == "" {
		http.Error(w, helpers.Error("name or description must be set"), 400)
		return
	}

	if input.Name != "" {
		template.Name = input.Name
	}

	if input.Description != "" {
		template.Description = zero.StringFrom(input.Description)
	}

	err = h.Datastore.SaveTemplate(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	idTemplate := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idTemplate, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteTemplateById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

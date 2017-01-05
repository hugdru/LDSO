package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"strconv"
)

func (h *Handler) templatesRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getTemplates))
	router.Get("/current", decorators.ReplyJson(h.getCurrentTemplate))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createTemplate)))
	router.Route("/:idt", h.templateRoutes)
}

func (h *Handler) templateRoutes(router chi.Router) {
	router.Use(h.templateContext)
	router.Get("/", decorators.ReplyJson(h.getTemplate))
	router.Post("/close", decorators.OnlySuperadmins(decorators.ReplyJson(h.closeTemplate)))
	router.Get("/used", decorators.ReplyJson(h.checkTemplateUsed))
	router.Get("/maingroups", decorators.ReplyJson(h.getTemplateMaingroups))
	router.Get("/subgroups", decorators.ReplyJson(h.getTemplateSubgroups))
	router.Get("/criteria", decorators.ReplyJson(h.getTemplateCriteria))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateTemplate)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteTemplate)))
}

func (h *Handler) templateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idTemplateStr := chi.URLParam(r, "idt")
		idTemplate, err := helpers.ParseInt64(idTemplateStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		template, err := h.Datastore.GetTemplateById(idTemplate)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "template", template)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

	templates, err := h.Datastore.GetTemplatesWithMaingroups(limit, offset, filter)
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

func (h *Handler) getCurrentTemplate(w http.ResponseWriter, r *http.Request) {
	template, err := h.Datastore.GetCurrentTemplate()
	if err != sql.ErrNoRows && err != nil {
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

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		input.Name = r.PostFormValue("name")
		input.Description = r.PostFormValue("description")
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

	template := datastore.NewTemplate(false)
	err := template.MustSet(input.Name)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = template.OptionalSetIfNotEmptyOrNil(input.Description)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	template.CreatedDate = helpers.TheTime()

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

func (h *Handler) getTemplate(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	templateWithMaingroups, err := h.Datastore.GetTemplateWithMaingroups(template.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	templateWithMaingroupsSlice, err := json.Marshal(templateWithMaingroups)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(templateWithMaingroupsSlice)
}

func (h *Handler) checkTemplateUsed(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	used, err := h.Datastore.CheckTemplateUsed(template.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	w.Write([]byte(helpers.FastConcat(`{"used":`, strconv.FormatBool(used), `}`)))
}

func (h *Handler) closeTemplate(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	var input struct {
		Close bool `json:"close"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.Close, err = strconv.ParseBool(r.PostFormValue("close"))
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

	template.Closed = input.Close
	if input.Close {
		template.ClosedDate = zero.TimeFrom(helpers.TheTime())
	} else {
		template.ClosedDate = zero.Time{}
	}
	err := h.Datastore.UpdateTemplate(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getTemplateMaingroups(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	templateMaingroups, err := h.Datastore.GetMaingroupsByTemplateId(template.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	templateMaingroupsSlice, err := json.Marshal(templateMaingroups)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(templateMaingroupsSlice)
}

func (h *Handler) getTemplateSubgroups(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	templateSubgroups, err := h.Datastore.GetSubgroupsByTemplateId(template.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	templateSubgroupsSlice, err := json.Marshal(templateSubgroups)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(templateSubgroupsSlice)
}

func (h *Handler) getTemplateCriteria(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	templateCriteria, err := h.Datastore.GetCriteriaByTemplateIdWithCriterionAccessibility(template.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	templateCriteriaSlice, err := json.Marshal(templateCriteria)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(templateCriteriaSlice)
}

func (h *Handler) updateTemplate(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		input.Name = r.PostFormValue("name")
		input.Description = r.PostFormValue("description")
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

	if input.Name == "" && input.Description == "" {
		http.Error(w, helpers.Error("[name] [description]"), 400)
		return
	}

	err := template.AllSetIfNotEmptyOrNil(input.Name, input.Description)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveTemplate(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteTemplate(w http.ResponseWriter, r *http.Request) {

	template := r.Context().Value("template").(*datastore.Template)

	err := h.Datastore.DeleteTemplate(template)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) templatesRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.GetTemplates))
	router.Get("/:id", helpers.ReplyJson(h.GetTemplate))
}

func (h *Handler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	var limit, offset int = 10, 0
	var err error

	limitString := r.FormValue("limit")
	offsetString := r.FormValue("offset")

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, helpers.Error("Invalid limit"), 400)
			return
		}
	}
	if offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			http.Error(w, helpers.Error("Invalid offset"), 400)
			return
		}
	}

	if limit <= 0 || offset < 0 || limit > 100 {
		http.Error(w, helpers.Error("0<limit<=100 && offset > 0"), 400)
		return
	}

	templates, err := h.Datastore.GetTemplates(limit, offset)
	if err != nil {
		http.Error(w, helpers.Error("Failed to grab Templates"), 500)
		return
	}
	templatesSlice, err := json.Marshal(templates)
	if err != nil {
		http.Error(w, helpers.Error("Failed converting Templates to JSON"), 500)
		return
	}
	w.Write(templatesSlice)
}

func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	//idTemplate := chi.URLParam(r, "id")
	//id, err := strconv.ParseInt(idTemplate, 10, 64)
	//if err != nil {
	//	http.Error(w, helpers.Error("Bad Template id"), 400)
	//	return
	//}
	//template, err := h.Datastore.GetTemplate(id)
	//if err != nil {
	//	http.Error(w, helpers.Error("Invalid Template id"), 400)
	//	return
	//}
	//templateSlice, err := json.Marshal(template)
	//if err != nil {
	//	http.Error(w, helpers.Error("Failed converting Templates to JSON"), 500)
	//	return
	//}
	//w.Write(templateSlice)
}

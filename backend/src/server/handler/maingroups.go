package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"
	"time"
	"server/handler/helpers/decorators"
)

func (h *Handler) maingroupsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getMaingroups))
	router.Post("/", decorators.RequestJson(decorators.ReplyJson(h.createMaingroup)))
	router.Get("/:id", decorators.ReplyJson(h.getMaingroup))
	router.Put("/:id", decorators.RequestJson(decorators.ReplyJson(h.updateMaingroup)))
	router.Delete("/:id", decorators.ReplyJson(h.deleteMaingroup))
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

func (h *Handler) getMaingroup(w http.ResponseWriter, r *http.Request) {
	idMaingroup := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idMaingroup, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	maingroup, err := h.Datastore.GetMaingroupById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	maingroupSlice, err := json.Marshal(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(maingroupSlice)
}

func (h *Handler) createMaingroup(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		IdTemplate int64
		Name       string
		Weight     int
	}
	input.Weight = -1

	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.IdTemplate == 0 || input.Name == "" || input.Weight == -1 {
		http.Error(w, helpers.Error("The maingroup must have a name idTemplate, a name and a weight"), 400)
		return
	}

	_, err = h.Datastore.GetTemplateById(input.IdTemplate)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	maingroup := datastore.NewMaingroup(false)
	maingroup.IdTemplate = input.IdTemplate
	maingroup.Name = input.Name
	maingroup.Weight = input.Weight
	maingroup.CreatedDate = time.Now().UTC()

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

func (h *Handler) updateMaingroup(w http.ResponseWriter, r *http.Request) {
	idMaingroup := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idMaingroup, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	maingroup, err := h.Datastore.GetMaingroupById(id)
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
		Name   string
		Weight int
	}
	input.Weight = -1

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" && input.Weight == -1 {
		http.Error(w, helpers.Error("name or weight must be set"), 400)
		return
	}

	if input.Name != "" {
		maingroup.Name = input.Name
	}

	if input.Weight != -1 {
		maingroup.Weight = input.Weight
	}

	err = h.Datastore.SaveMaingroup(maingroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteMaingroup(w http.ResponseWriter, r *http.Request) {
	idMaingroup := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idMaingroup, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteMaingroupById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

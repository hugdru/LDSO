package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"strconv"
)

func (h *Handler) subgroupsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getSubgroups))
	router.Post("/", decorators.RequestJson(decorators.ReplyJson(h.createSubgroup)))
	router.Get("/:id", decorators.ReplyJson(h.getSubgroup))
	router.Put("/:id", decorators.RequestJson(decorators.ReplyJson(h.updateSubgroup)))
	router.Delete("/:id", decorators.ReplyJson(h.deleteSubgroup))
}

func (h *Handler) getSubgroups(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"idMaingroup", "id_maingroup"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	subgroups, err := h.Datastore.GetSubgroups(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	subgroupsSlice, err := json.Marshal(subgroups)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(subgroupsSlice)
}

func (h *Handler) getSubgroup(w http.ResponseWriter, r *http.Request) {
	idSubgroup := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idSubgroup, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	subgroup, err := h.Datastore.GetSubgroupById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	subgroupSlice, err := json.Marshal(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(subgroupSlice)
}

func (h *Handler) createSubgroup(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		IdMaingroup int64
		Name        string
		Weight      int
	}
	input.Weight = -1

	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.IdMaingroup == 0 || input.Name == "" || input.Weight == -1 {
		http.Error(w, helpers.Error("The subgroup must have a name idMaingroup, a name and a weight"), 400)
		return
	}

	_, err = h.Datastore.GetMaingroupById(input.IdMaingroup)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	subgroup := datastore.NewSubgroup(false)
	subgroup.IdMaingroup = input.IdMaingroup
	subgroup.Name = input.Name
	subgroup.Weight = input.Weight
	subgroup.CreatedDate = helpers.TheTime()

	err = h.Datastore.SaveSubgroup(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	subgroupSlice, err := json.Marshal(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(subgroupSlice)
}

func (h *Handler) updateSubgroup(w http.ResponseWriter, r *http.Request) {
	idSubgroup := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idSubgroup, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	subgroup, err := h.Datastore.GetSubgroupById(id)
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
		subgroup.Name = input.Name
	}

	if input.Weight != -1 {
		subgroup.Weight = input.Weight
	}

	err = h.Datastore.SaveSubgroup(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteSubgroup(w http.ResponseWriter, r *http.Request) {
	idSubgroup := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idSubgroup, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteSubgroupById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

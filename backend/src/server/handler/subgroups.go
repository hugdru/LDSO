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

func (h *Handler) subgroupsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getSubgroups))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createSubgroup)))
	router.Route("/:ids", h.subgroupRoutes)

}

func (h *Handler) subgroupRoutes(router chi.Router) {
	router.Use(h.subgroupContext)
	router.Get("/", decorators.ReplyJson(h.getSubgroup))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateSubgroup)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteSubgroup)))
}

func (h *Handler) subgroupContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idSubgroupStr := chi.URLParam(r, "ids")
		idSubgroup, err := helpers.ParseInt64(idSubgroupStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		subgroup, err := h.Datastore.GetSubgroupById(idSubgroup)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "subgroup", subgroup)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func (h *Handler) createSubgroup(w http.ResponseWriter, r *http.Request) {

	var input struct {
		IdMaingroup int64  `json:"idMaingroup"`
		Name        string `json:"name"`
		Weight      int    `json:"weight"`
	}
	input.Weight = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdMaingroup, err = helpers.ParseInt64(r.PostFormValue("idMaingroup"))
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

	subgroup := datastore.NewSubgroup(false)
	err := subgroup.MustSet(input.IdMaingroup, input.Name, input.Weight)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
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

func (h *Handler) getSubgroup(w http.ResponseWriter, r *http.Request) {

	subgroup := r.Context().Value("subgroup").(*datastore.Subgroup)

	subgroupSlice, err := json.Marshal(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(subgroupSlice)
}

func (h *Handler) updateSubgroup(w http.ResponseWriter, r *http.Request) {

	subgroup := r.Context().Value("subgroup").(*datastore.Subgroup)

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

	err := subgroup.UpdateSetIfNotEmptyOrNil(input.Name, input.Weight)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveSubgroup(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteSubgroup(w http.ResponseWriter, r *http.Request) {

	subgroup := r.Context().Value("subgroup").(*datastore.Subgroup)

	err := h.Datastore.DeleteSubgroup(subgroup)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

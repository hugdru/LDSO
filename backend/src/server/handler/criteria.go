package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"
	"time"
)

func (h *Handler) criteriaRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getCriteria))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createCriterion)))
	router.Get("/:id", helpers.ReplyJson(h.getCriterion))
	router.Put("/:id", helpers.RequestJson(helpers.ReplyJson(h.updateCriterion)))
	router.Delete("/:id", helpers.ReplyJson(h.deleteCriterion))
}

func (h *Handler) getCriteria(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"idSubgroup", "id_subgroup"},
		[]string{"idLegislation", "id_legislation"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	criteria, err := h.Datastore.GetCriteria(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	criteriaSlice, err := json.Marshal(criteria)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(criteriaSlice)
}

func (h *Handler) getCriterion(w http.ResponseWriter, r *http.Request) {
	idCriterion := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idCriterion, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	criterion, err := h.Datastore.GetCriterionById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	criterionSlice, err := json.Marshal(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(criterionSlice)
}

func (h *Handler) createCriterion(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		IdSubgroup    int64
		IdLegislation int64
		Name          string
		Weight        int
	}
	input.Weight = -1

	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.IdSubgroup == 0 || input.Name == "" || input.Weight == -1 {
		http.Error(w, helpers.Error("idSubgroup name weight [idLegislation]"), 400)
		return
	}

	if input.IdLegislation != 0 {
		_, err = h.Datastore.GetLegislationById(input.IdLegislation)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

	criterion := datastore.NewCriterion(false)
	criterion.IdSubgroup = input.IdSubgroup
	criterion.IdLegislation = zero.IntFrom(input.IdLegislation)
	criterion.Name = input.Name
	criterion.Weight = input.Weight
	criterion.CreatedDate = time.Now().UTC()

	err = h.Datastore.SaveCriterion(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	criterionSlice, err := json.Marshal(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(criterionSlice)
}

func (h *Handler) updateCriterion(w http.ResponseWriter, r *http.Request) {
	idCriterion := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idCriterion, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	criterion, err := h.Datastore.GetCriterionById(id)
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
		IdLegislation int64
		Name          string
		Weight        int
	}
	input.Weight = -1

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Name == "" && input.Weight == -1 && input.IdLegislation == 0 {
		http.Error(w, helpers.Error("name or weight or idLegislation must be set"), 400)
		return
	}

	if input.Name != "" {
		criterion.Name = input.Name
	}

	if input.Weight != -1 {
		criterion.Weight = input.Weight
	}

	if input.IdLegislation != 0 {
		_, err = h.Datastore.GetLegislationById(input.IdLegislation)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		criterion.IdLegislation = zero.IntFrom(input.IdLegislation)
	}

	err = h.Datastore.SaveCriterion(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteCriterion(w http.ResponseWriter, r *http.Request) {
	idCriterion := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idCriterion, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteCriterionById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

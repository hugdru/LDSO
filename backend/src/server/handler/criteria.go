package handler

import (
	"context"
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

	router.Route("/:id/accessibilities", h.criteriaAccessibilitiesSubroutes)
}

func (h *Handler) criteriaAccessibilitiesSubroutes(router chi.Router) {
	router.Use(h.criteriaAccessibilitiesContext)
	router.Get("/", helpers.ReplyJson(h.getCriterionAccessibilities))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createCriterionAccessibility)))
	router.Delete("/", helpers.ReplyJson(h.deleteCriterionAccessibilities))
	router.Get("/:ida", helpers.ReplyJson(h.getCriterionAccessibility))
	router.Put("/:ida", helpers.RequestJson(helpers.ReplyJson(h.updateCriterionAccessibility)))
	router.Delete("/:ida", helpers.ReplyJson(h.deleteCriterionAccessibility))
}

func (h *Handler) criteriaAccessibilitiesContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idCriterionStr := chi.URLParam(r, "id")
		idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		criterion, err := h.Datastore.GetCriterionById(idCriterion)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		ctx := context.WithValue(r.Context(), "criterion", criterion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func (h *Handler) getCriterionAccessibilities(w http.ResponseWriter, r *http.Request) {
	criterion := r.Context().Value("criterion").(*datastore.Criterion)
	accessibilities, err := h.Datastore.GetAccessibilitiesByCriterionId(criterion.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	accessibilitiesSlice, err := json.Marshal(accessibilities)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(accessibilitiesSlice)
}

func (h *Handler) deleteCriterionAccessibilities(w http.ResponseWriter, r *http.Request) {
	criterion := r.Context().Value("criterion").(*datastore.Criterion)
	err := h.Datastore.DeleteAccessibilitiesByCriterionId(criterion.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getCriterionAccessibility(w http.ResponseWriter, r *http.Request) {
	idAccessibilityStr := chi.URLParam(r, "ida")
	idAccessibility, err := strconv.ParseInt(idAccessibilityStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	accessibility, err := h.Datastore.GetAccessibilityByIds(criterion.Id, idAccessibility)
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

func (h *Handler) createCriterionAccessibility(w http.ResponseWriter, r *http.Request) {
	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	var input struct {
		IdAccessibility int64
		Weight          int
	}
	input.Weight = -1

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Weight == -1 || input.IdAccessibility == 0 {
		http.Error(w, helpers.Error("weight and idAccessibility must be set"), 400)
		return
	}

	err = h.Datastore.InsertAccessibilityByCriterionId(criterion.Id, input.IdAccessibility, input.Weight)

}

func (h *Handler) updateCriterionAccessibility(w http.ResponseWriter, r *http.Request) {
	idAccessibilityStr := chi.URLParam(r, "ida")
	idAccessibility, err := strconv.ParseInt(idAccessibilityStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	accessibility, err := h.Datastore.GetAccessibilityByIds(criterion.Id, idAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	d := json.NewDecoder(r.Body)
	if d == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Weight int
	}
	input.Weight = -1

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Weight == -1 {
		http.Error(w, helpers.Error("weight must be set"), 400)
		return
	}

	accessibility.Weight = input.Weight

	err = h.Datastore.SaveAccessibilityByCriterionIdAccessibility(criterion.Id, accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteCriterionAccessibility(w http.ResponseWriter, r *http.Request) {
	idAccessibilityStr := chi.URLParam(r, "ida")
	idAccessibility, err := strconv.ParseInt(idAccessibilityStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	err = h.Datastore.DeleteAccessibilitiesByIds(criterion.Id, idAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}
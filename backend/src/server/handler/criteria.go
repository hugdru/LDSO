package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
)

func (h *Handler) criteriaRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getCriteria))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createCriterion)))
	router.Route("/:idc", h.criterionSubroutes)
}

func (h *Handler) criterionSubroutes(router chi.Router) {
	router.Use(h.criterionContext)
	router.Get("/", decorators.ReplyJson(h.getCriterion))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateCriterion)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteCriterion)))
	router.Route("/accessibilities", h.criterionAccessibilitiesSubroutes)
}

func (h *Handler) criterionAccessibilitiesSubroutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getCriterionAccessibilities))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createCriterionAccessibility)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteCriterionAccessibilities)))
	router.Route("/:ida", h.criterionAccessibilitySubroutes)
}

func (h *Handler) criterionAccessibilitySubroutes(router chi.Router) {
	router.Use(h.criterionAccessibilityContext)
	router.Get("/", decorators.ReplyJson(h.getCriterionAccessibility))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateCriterionAccessibility)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteCriterionAccessibility)))
}

func (h *Handler) criterionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idCriterionStr := chi.URLParam(r, "idc")
		idCriterion, err := helpers.ParseInt64(idCriterionStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		criterion, err := h.Datastore.GetCriterionByIdWithLegislation(idCriterion)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		ctx := context.WithValue(r.Context(), "criterion", criterion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) criterionAccessibilityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAccessibilityStr := chi.URLParam(r, "ida")
		idAccessibility, err := helpers.ParseInt64(idAccessibilityStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		ctx := context.WithValue(r.Context(), "idAccessibility", idAccessibility)
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

	criteria, err := h.Datastore.GetCriteriaWithLegislation(limit, offset, filter)
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

type inputCriterionLegislation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func (h *Handler) createCriterion(w http.ResponseWriter, r *http.Request) {

	var input struct {
		IdSubgroup    int64  `json:"idSubgroup"`
		IdLegislation int64  `json:"idLegislation"`
		Name          string `json:"name"`
		Weight        int    `json:"weight"`
		Legislation   string `json:"legislation"`
		//Legislation   *inputCriterionLegislation `json:"legislation"`
	}
	input.Weight = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdSubgroup, err = helpers.ParseInt64(r.PostFormValue("idSubgroup"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.IdLegislation, _ = helpers.ParseInt64(r.PostFormValue("idLegislation"))
		input.Name = r.PostFormValue("name")
		input.Weight, err = helpers.ParseInt(r.PostFormValue("weight"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.Legislation = r.PostFormValue("legislation")
		//input.Legislation.Name = r.PostFormValue("legislation.name")
		//input.Legislation.Description = r.PostFormValue("legislation.description")
		//input.Legislation.Url = r.PostFormValue("legislation.url")
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

	newCriterion := datastore.NewCriterion(false)
	err := newCriterion.MustSet(input.IdSubgroup, input.Name, input.Weight)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = newCriterion.OptionalSetIfNotEmptyOrNil(input.IdLegislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	newCriterion.CreatedDate = helpers.TheTime()

	if input.Legislation != "" {
		resultIdLegislation, err := insertOrFetchLegislation(h.Datastore, input.Legislation)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
		newCriterion.IdLegislation = zero.IntFrom(resultIdLegislation)
	} else {
		newCriterion.IdLegislation = zero.IntFrom(0)
	}

	/*
		if input.IdLegislation == 0 && input.Legislation != nil {
				resultIdLegislation, err := insertOrFetchLegislation(h.Datastore, input.Legislation)
				if err != nil {
					http.Error(w, helpers.Error(err.Error()), 500)
					return
				}
				newCriterion.IdLegislation = zero.IntFrom(resultIdLegislation)
			}
	*/

	err = h.Datastore.SaveCriterion(newCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	accessibilities, err := h.Datastore.GetAccessibilities(nil)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	for _, accessibility := range accessibilities {
		err = h.Datastore.InsertCriterionAccessibilityByIds(newCriterion.Id, accessibility.Id, 0)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	}

	criterionSlice, err := json.Marshal(newCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(criterionSlice)
}

func (h *Handler) getCriterion(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	criterionSlice, err := json.Marshal(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(criterionSlice)
}

func (h *Handler) updateCriterion(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	var input struct {
		IdLegislation int64  `json:"idLegislation"`
		Name          string `json:"name"`
		Weight        int    `json:"weight"`
		Legislation   string `json:"legislation"`
		//Legislation   *inputCriterionLegislation `json:"legislation"`
	}
	input.Weight = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdLegislation, _ = helpers.ParseInt64(r.PostFormValue("idLegislation"))
		input.Name = r.PostFormValue("name")
		input.Weight, _ = helpers.ParseInt(r.PostFormValue("weight"))
		if err != nil {
			input.Weight = -1
		}
		input.Legislation = r.PostFormValue("legislation")
		//input.Legislation.Name = r.PostFormValue("legislation.name")
		//input.Legislation.Description = r.PostFormValue("legislation.description")
		//input.Legislation.Url = r.PostFormValue("legislation.url")
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

	//if input.Name == "" && input.Weight == -1 && input.IdLegislation == 0 && input.Legislation == nil {
	//	http.Error(w, helpers.Error("[name] [weight] [idLegislation] [legislation]"), 400)
	//	return
	//}

	if input.Name == "" && input.Weight == -1 && input.IdLegislation == 0 && input.Legislation == "" {
		http.Error(w, helpers.Error("[name] [weight] [idLegislation] [legislation]"), 400)
		return
	}

	err := criterion.UpdateSetIfNotEmptyOrNil(input.Name, input.Weight, input.IdLegislation)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Legislation != "" {
		resultIdLegislation, err := insertOrFetchLegislation(h.Datastore, input.Legislation)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
		criterion.IdLegislation = zero.IntFrom(resultIdLegislation)
	} else {
		criterion.IdLegislation = zero.IntFrom(0)
	}

	/*
		if input.IdLegislation == 0 && input.Legislation != nil {
			resultIdLegislation, err := insertOrFetchLegislation(h.Datastore, input.Legislation)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 500)
				return
			}
			criterion.IdLegislation = zero.IntFrom(resultIdLegislation)
		}
	*/

	err = h.Datastore.SaveCriterion(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteCriterion(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	err := h.Datastore.DeleteCriterion(criterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getCriterionAccessibilities(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	criterionAccessibilities, err := h.Datastore.GetCriterionAccessibilitiesByCriterionId(criterion.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	criterionAccessibilitiesSlice, err := json.Marshal(criterionAccessibilities)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(criterionAccessibilitiesSlice)
}

func (h *Handler) createCriterionAccessibility(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	var input struct {
		IdAccessibility int64 `json:"idAccessibility"`
		Weight          int   `json:"weight"`
	}
	input.Weight = -1

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		var err error
		input.IdAccessibility, err = helpers.ParseInt64(r.PostFormValue("idAccessibility"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
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

	if input.IdAccessibility == 0 || input.Weight == -1 {
		http.Error(w, helpers.Error("idAccessibility weight"), 400)
		return
	}

	err := h.Datastore.InsertCriterionAccessibilityByIds(criterion.Id, input.IdAccessibility, input.Weight)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

}

func (h *Handler) deleteCriterionAccessibilities(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	err := h.Datastore.DeleteCriterionAccessibilitiesByCriterionId(criterion.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getCriterionAccessibility(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	idAccessibility := r.Context().Value("idAccessibility").(int64)

	criterionAccessibility, err := h.Datastore.GetCriterionAccessibilityByIds(criterion.Id, idAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	criterionAccessibilitySlice, err := json.Marshal(criterionAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(criterionAccessibilitySlice)
}

func (h *Handler) updateCriterionAccessibility(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	idAccessibility := r.Context().Value("idAccessibility").(int64)

	criterionAccessibility, err := h.Datastore.GetCriterionAccessibilityByIds(criterion.Id, idAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	var input struct {
		Weight int `json:"weight"`
	}
	input.Weight = -1

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		var err error
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

	if input.Weight == -1 {
		http.Error(w, helpers.Error("weight must be set"), 400)
		return
	}

	criterionAccessibility.Weight = input.Weight

	err = h.Datastore.SaveCriterionAccessibility(criterionAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteCriterionAccessibility(w http.ResponseWriter, r *http.Request) {

	criterion := r.Context().Value("criterion").(*datastore.Criterion)

	idAccessibility := r.Context().Value("idAccessibility").(int64)

	err := h.Datastore.DeleteCriterionAccessibilityByIds(criterion.Id, idAccessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func insertOrFetchLegislation(d *datastore.Datastore,
	inputLegislationName string) (int64, error) {
	if inputLegislationName == "" {
		return 0, errors.New("At least a name must be specified")
	}
	legislation, err := d.GetLegislationByName(inputLegislationName)
	if err == sql.ErrNoRows {
		legislation = datastore.NewLegislation(true)
		legislation.Name = inputLegislationName
		err = d.InsertLegislation(legislation)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	return legislation.Id, nil
}

/*
func insertOrFetchLegislation(d *datastore.Datastore, inputLegislation *inputCriterionLegislation) (int64, error) {
	if inputLegislation.Name == "" {
		return 0, errors.New("At least a name must be specified")
	}

	legislation, err := d.GetLegislationByName(inputLegislation.Name)

	if err == sql.ErrNoRows {
		legislation = datastore.NewLegislation(false)
		err = legislation.AllSetIfNotEmptyOrNil(inputLegislation.Name, inputLegislation.Description, inputLegislation.Url)
		if err != nil {
			return 0, err
		}
		err = d.InsertLegislation(legislation)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	return legislation.Id, nil
}
*/

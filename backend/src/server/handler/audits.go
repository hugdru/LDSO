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

func (h *Handler) auditsRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getAudits))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createAudit)))
	router.Get("/:id", helpers.ReplyJson(h.getAudit))
	router.Put("/:id", helpers.RequestJson(helpers.ReplyJson(h.updateAudit)))
	router.Delete("/:id", helpers.ReplyJson(h.deleteAudit))

	router.Get("/:id/subgroups", helpers.ReplyJson(h.getAuditSubgroups))
	router.Post("/:id/subgroups", helpers.RequestJson(helpers.ReplyJson(h.createAuditSubgroups)))
	router.Delete("/:id/subgroups", helpers.ReplyJson(h.deleteAuditSubgroups))

	router.Get("/:id/criteria", helpers.ReplyJson(h.getAuditCriteria))
	router.Get("/:id/criteria/:idc", helpers.ReplyJson(h.getAuditCriterion))
	router.Post("/:id/criteria/:idc", helpers.RequestJson(helpers.ReplyJson(h.createAuditCriterion)))
	router.Put("/:id/criteria/:idc", helpers.RequestJson(helpers.ReplyJson(h.updateAuditCriterion)))
	router.Delete("/:id/criteria", helpers.ReplyJson(h.deleteAuditCriteria))
	router.Delete("/:id/criteria/:idc", helpers.ReplyJson(h.deleteAuditCriterion))
}

func (h *Handler) getAudits(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"idProperty", "id_property"},
		[]string{"idAuditor", "id_auditor"},
		[]string{"idTemplate", "id_template"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	audits, err := h.Datastore.GetAudits(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditsSlice, err := json.Marshal(audits)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditsSlice)
}

func (h *Handler) createAudit(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		IdProperty int64
		IdAuditor  int64
		IdTemplate int64
	}

	err := decoder.Decode(&input)

	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.IdProperty == 0 || input.IdAuditor == 0 || input.IdTemplate == 0 {
		http.Error(w, helpers.Error("The audits must have idProperty, idAuditor and idTemplate"), 400)
		return
	}

	_, err = h.Datastore.GetPropertyById(input.IdProperty)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_, err = h.Datastore.GetAuditorByIdWithForeign(input.IdAuditor)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_, err = h.Datastore.GetTemplateById(input.IdTemplate)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	audit := datastore.NewAudit(false)
	audit.IdProperty = input.IdProperty
	audit.IdAuditor = input.IdAuditor
	audit.IdTemplate = input.IdTemplate
	audit.CreatedDate = time.Now().UTC()

	err = h.Datastore.SaveAudit(audit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditSlice, err := json.Marshal(audit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditSlice)
}

func (h *Handler) getAudit(w http.ResponseWriter, r *http.Request) {
	idAudit := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAudit, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	audit, err := h.Datastore.GetAuditById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	auditSlice, err := json.Marshal(audit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditSlice)
}
func (h *Handler) updateAudit(w http.ResponseWriter, r *http.Request) {
	idAudit := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAudit, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	audit, err := h.Datastore.GetAuditById(id)
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
		IdAuditor   int64
		Rating      int64
		Observation string
	}
	input.Rating = -1

	err = d.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.IdAuditor == 0 && input.Rating == -1 && input.Observation == "" {
		http.Error(w, helpers.Error("At least one of IdAuditor, Rating or Observation"), 400)
		return
	}

	if input.IdAuditor != 0 {
		audit.IdAuditor = input.IdAuditor
	}

	if input.Rating != -1 {
		audit.Rating = zero.IntFrom(input.Rating)
	}

	if input.Observation != "" {
		audit.Observation = zero.StringFrom(input.Observation)
	}

	err = h.Datastore.SaveAudit(audit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteAudit(w http.ResponseWriter, r *http.Request) {
	idAudit := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAudit, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getAuditSubgroups(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"idAudit", "id_audit"},
		[]string{"idSubgroup", "id_subgroup"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	subgroups, err := h.Datastore.GetAuditSubgroupsByIdAudit(idAudit, filter)
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

func (h *Handler) createAuditSubgroups(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	audit, err := h.Datastore.GetAuditById(idAudit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	idTemplate := audit.IdTemplate

	decoder := json.NewDecoder(r.Body)

	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input []int64

	err = decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if len(input) == 0 {
		http.Error(w, helpers.Error("At least one idSubgroup"), 400)
		return
	}

	err = h.Datastore.SaveAuditSubgroup(idAudit, idTemplate, input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteAuditSubgroups(w http.ResponseWriter, r *http.Request) {
	idAudit := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAudit, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditSubgroupsByIdAudit(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getAuditCriteria(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := make(map[string]interface{})
	filter["id_audit"] = idAudit

	auditCriteria, err := h.Datastore.GetAuditCriteria(filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditCriteriaSlice, err := json.Marshal(auditCriteria)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditCriteriaSlice)
}

func (h *Handler) getAuditCriterion(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterion, err := h.Datastore.GetAuditCriterionByIds(idAudit, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditCriterionSlice, err := json.Marshal(auditCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditCriterionSlice)
}

func (h *Handler) createAuditCriterion(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	decoder := json.NewDecoder(r.Body)

	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Value       int64
		Observation string
	}

	err = decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterion := datastore.NewAuditCriterion(false)
	auditCriterion.IdAudit = idAudit
	auditCriterion.IdCriterion = idCriterion
	auditCriterion.Value = zero.IntFrom(input.Value)
	auditCriterion.Observation = zero.StringFrom(input.Observation)

	err = h.Datastore.SaveAuditCriterion(auditCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditCriterionSlice, err := json.Marshal(auditCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditCriterionSlice)
}

func (h *Handler) updateAuditCriterion(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterion, err := h.Datastore.GetAuditCriterionByIds(idAudit, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		Value       int64
		Observation string
	}
	input.Value = -1

	err = decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Value == -1 && input.Observation == "" {
		http.Error(w, helpers.Error("At least one of value, observation"), 400)
		return
	}

	if input.Value != -1 {
		auditCriterion.Value = zero.IntFrom(input.Value)
	}
	if input.Observation != "" {
		auditCriterion.Observation = zero.StringFrom(input.Observation)
	}

	err = h.Datastore.SaveAuditCriterion(auditCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditCriterionSlice, err := json.Marshal(auditCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditCriterionSlice)
}

func (h *Handler) deleteAuditCriteria(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditCriterionByIdAudit(idAudit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) deleteAuditCriterion(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "id")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditCriterionByIds(idAudit, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

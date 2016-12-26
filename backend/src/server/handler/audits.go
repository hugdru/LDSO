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

func (h *Handler) auditsRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getAudits))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createAudit)))

	router.Get("/:ida", helpers.ReplyJson(h.getAudit))
	router.Put("/:ida", helpers.RequestJson(helpers.ReplyJson(h.updateAudit)))
	router.Delete("/:ida", helpers.ReplyJson(h.deleteAudit))

	router.Get("/:ida/subgroups", helpers.ReplyJson(h.getAuditSubgroups))
	router.Post("/:ida/subgroups", helpers.RequestJson(helpers.ReplyJson(h.createAuditSubgroups)))
	router.Delete("/:ida/subgroups", helpers.ReplyJson(h.deleteAuditSubgroups))

	router.Route("/:ida/criteria", h.auditsCriteriaSubroutes)
}

func (h *Handler) auditsCriteriaSubroutes(router chi.Router) {
	router.Use(h.auditsCriteriaContext)

	router.Get("/", helpers.ReplyJson(h.getAuditCriteria))
	router.Delete("/", helpers.ReplyJson(h.deleteAuditCriteria))
	router.Get("/:idc", helpers.ReplyJson(h.getAuditCriterion))
	router.Post("/:idc", helpers.RequestJson(helpers.ReplyJson(h.createAuditCriterion)))
	router.Put("/:idc", helpers.RequestJson(helpers.ReplyJson(h.updateAuditCriterion)))
	router.Delete("/:idc", helpers.ReplyJson(h.deleteAuditCriterion))

	router.Route("/:idc/remarks", h.auditsCriteriaRemarksSubroutes)
}

func (h *Handler) auditsCriteriaRemarksSubroutes(router chi.Router) {
	router.Post("/", h.createCriterionRemark)
	router.Get("/:idr/image", h.getCriterionRemarkImage)
	router.Get("/:idr", h.getCriterionRemark)
	router.Get("/", h.getCriterionRemarks)
}

func (h *Handler) auditsCriteriaContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAuditStr := chi.URLParam(r, "ida")
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
		ctx := context.WithValue(r.Context(), "audit", audit)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
	idAuditStr := chi.URLParam(r, "ida")
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
	auditSlice, err := json.Marshal(audit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditSlice)
}
func (h *Handler) updateAudit(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "ida")
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
	idAuditStr := chi.URLParam(r, "ida")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditById(idAudit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getAuditSubgroups(w http.ResponseWriter, r *http.Request) {
	idAuditStr := chi.URLParam(r, "ida")
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
	idAuditStr := chi.URLParam(r, "ida")
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
	idAuditStr := chi.URLParam(r, "ida")
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditSubgroupsByIdAudit(idAudit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getAuditCriteria(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	filter := make(map[string]interface{})
	filter["id_audit"] = audit.Id

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

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterion, err := h.Datastore.GetAuditCriterionByIds(audit.Id, idCriterion)
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

	audit := r.Context().Value("audit").(*datastore.Audit)

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
		Value int64
	}

	err = decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterion := datastore.NewAuditCriterion(false)
	auditCriterion.IdAudit = audit.Id
	auditCriterion.IdCriterion = idCriterion
	auditCriterion.Value = zero.IntFrom(input.Value)

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

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterion, err := h.Datastore.GetAuditCriterionByIds(audit.Id, idCriterion)
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
		Value int64
	}
	input.Value = -1

	err = decoder.Decode(&input)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.Value == -1 {
		http.Error(w, helpers.Error("At least one of value, observation"), 400)
		return
	}

	if input.Value != -1 {
		auditCriterion.Value = zero.IntFrom(input.Value)
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

	audit := r.Context().Value("audit").(*datastore.Audit)

	err := h.Datastore.DeleteAuditCriterionByIdAudit(audit.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) deleteAuditCriterion(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.DeleteAuditCriterionByIds(audit.Id, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) createCriterionRemark(w http.ResponseWriter, r *http.Request) {

	contentLength := helpers.GetContentLength(r.Header.Get("Content-Length"))
	if contentLength == -1 {
		http.Error(w, helpers.Error("Invalid Content-Length header value"), http.StatusBadRequest)
		return
	}

	if contentLength > helpers.MaxMultipartSize {
		http.Error(w, helpers.Error("Data too big"), http.StatusBadRequest)
		return
	}

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterionStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		Observation   string
		imageMimetype string
		imageBytes    []byte
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		input.Observation = r.PostFormValue("Observation")
		input.imageBytes, input.imageMimetype, err = helpers.ReadImage(r, "image", helpers.MaxImageFileSize)
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	default:
		http.Error(w, helpers.Error("Content-type not supported"), 415)
		return
	}

	if input.Observation == "" && input.imageBytes == nil {
		http.Error(w, helpers.Error("observation or imageBytes"), 400)
		return
	}

	remark := datastore.NewRemark(false)
	remark.IdAudit = audit.Id
	remark.IdCriterion = idCriterion
	remark.Observation = zero.StringFrom(input.Observation)
	remark.Image = input.imageBytes
	remark.ImageMimeType = zero.StringFrom(input.imageMimetype)

	err = h.Datastore.InsertRemark(remark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"id":"` + helpers.Int64ToString(remark.Id) + "}"))
}

func (h *Handler) getCriterionRemark(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriteriaStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriteriaStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idRemarkStr := chi.URLParam(r, "idr")
	idRemark, err := strconv.ParseInt(idRemarkStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterionRemark, err := h.Datastore.GetRemarkByIdsAuditCriterionRemark(audit.Id, idCriterion, idRemark)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), http.StatusBadRequest)
		return
	}

	auditCriterionRemarkSlice, err := json.Marshal(auditCriterionRemark)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditCriterionRemarkSlice)
}

func (h *Handler) getCriterionRemarkImage(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriteriaStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriteriaStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idRemarkStr := chi.URLParam(r, "idr")
	idRemark, err := strconv.ParseInt(idRemarkStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditCriterionRemark, err := h.Datastore.GetRemarkByIdsAuditCriterionRemark(audit.Id, idCriterion, idRemark)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), http.StatusBadRequest)
		return
	}

	if auditCriterionRemark.ImageMimeType.Valid && auditCriterionRemark.Image != nil {
		w.Header().Set("Content-type", auditCriterionRemark.ImageMimeType.String)
		w.Header().Set("Content-Length", strconv.Itoa(len(auditCriterionRemark.Image)))
		w.Write(auditCriterionRemark.Image)
	} else {
		http.Error(w, helpers.Error("no image"), http.StatusBadRequest)
		return
	}
}

func (h *Handler) getCriterionRemarks(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriteriaStr := chi.URLParam(r, "idc")
	idCriterion, err := strconv.ParseInt(idCriteriaStr, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditsCriteriaRemarks, err := h.Datastore.GetRemarksByAuditCriterionIds(audit.Id, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditsCriteriaRemarksSlice, err := json.Marshal(auditsCriteriaRemarks)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	w.Write(auditsCriteriaRemarksSlice)
}

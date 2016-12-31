package handler

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/guregu/null.v3/zero"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"server/handler/sessionData"
)

func (h *Handler) auditsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getAudits))
	router.Post("/", decorators.OnlySuperadminsOrLocaladmins(decorators.ReplyJson(h.createAudit)))
	router.Route("/:ida", h.auditRoutes)
}

func (h *Handler) auditRoutes(router chi.Router) {
	router.Use(h.auditContext)
	router.Get("/", decorators.ReplyJson(h.getAudit))
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.updateAudit)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteAudit)))
	router.Route("/subgroups", h.auditSubgroupsRoutes)
	router.Route("/criteria", h.auditCriteriaSubroutes)
}

func (h *Handler) auditSubgroupsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getAuditSubgroups))
	router.Post("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.createAuditSubgroups)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteAuditSubgroups)))
}

func (h *Handler) auditCriteriaSubroutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getAuditCriteria))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteAuditCriteria)))
	router.Route("/:idc", h.auditCriterionSubroutes)
}

func (h *Handler) auditCriterionSubroutes(router chi.Router) {
	router.Use(h.auditCriterionContext)
	router.Get("/", decorators.ReplyJson(h.getAuditCriterion))
	router.Post("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.createAuditCriterion)))
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.updateAuditCriterion)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteAuditCriterion)))
	router.Route("/remarks", h.auditCriterionRemarksSubroutes)
}

func (h *Handler) auditCriterionRemarksSubroutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getCriterionRemarks))
	router.Post("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.createCriterionRemark)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteCriterionRemarks)))
	router.Route("/:idr", h.auditCriterionRemarkSubroutes)
}

func (h *Handler) auditCriterionRemarkSubroutes(router chi.Router) {
	router.Use(h.auditCriterionRemarkContext)
	router.Get("/", decorators.ReplyJson(h.getCriterionRemark))
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.updateCriterionRemark)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrAuditors(decorators.ReplyJson(h.deleteCriterionRemark)))
	router.Get("/image/:hash", h.getCriterionRemarkImage)
}

func (h *Handler) auditContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAuditStr := chi.URLParam(r, "ida")
		idAudit, err := helpers.ParseInt64(idAuditStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		audit, err := h.Datastore.GetAuditById(idAudit)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		esd, err := sessionData.GetSessionData(r)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		if r.Method != http.MethodGet {
			if !esd.IsSuperadmin() && !esd.IsLocaladmin() && esd.Id != audit.IdAuditor {
				http.Error(w, helpers.Error("Not the assigned auditor"), http.StatusForbidden)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "audit", audit)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) auditCriterionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idCriterionStr := chi.URLParam(r, "idc")
		idCriterion, err := helpers.ParseInt64(idCriterionStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		ctx := context.WithValue(r.Context(), "idCriterion", idCriterion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) auditCriterionRemarkContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idRemarkStr := chi.URLParam(r, "idr")
		idRemark, err := helpers.ParseInt64(idRemarkStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		ctx := context.WithValue(r.Context(), "idRemark", idRemark)
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

	var input struct {
		IdProperty int64 `json:"idProperty"`
		IdAuditor  int64 `json:"idAuditor"`
		IdTemplate int64 `json:"idTemplate"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdProperty, err = helpers.ParseInt64(r.PostFormValue("idProperty"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.IdAuditor, err = helpers.ParseInt64(r.PostFormValue("idAuditor"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.IdTemplate, err = helpers.ParseInt64(r.PostFormValue("idTemplate"))
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

	audit := datastore.NewAudit(false)
	err := audit.MustSet(input.IdProperty, input.IdAuditor, input.IdTemplate)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	audit.CreatedDate = helpers.TheTime()

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

	audit := r.Context().Value("audit").(*datastore.Audit)

	auditSlice, err := json.Marshal(audit)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(auditSlice)
}

func (h *Handler) updateAudit(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	esd, err := sessionData.GetSessionData(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if audit.IdAuditor != esd.Id && !esd.IsSuperadmin() && !esd.IsLocaladmin() {
		http.Error(w, helpers.Error("Not the owner of the property"), http.StatusForbidden)
		return
	}

	var input struct {
		IdAuditor   int64
		Rating      int64
		Observation string
	}
	input.Rating = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdAuditor, err = helpers.ParseInt64(r.PostFormValue("idAuditor"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.Rating, err = helpers.ParseInt64(r.PostFormValue("rating"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.Observation = r.PostFormValue("observation")
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

	if input.IdAuditor == 0 && input.Rating == -1 && input.Observation == "" {
		http.Error(w, helpers.Error("[idAuditor] [Rating] [Observation]"), 400)
		return
	}

	if input.IdAuditor != 0 {
		if esd.IsSuperadmin() || esd.IsLocaladmin() {
			audit.IdAuditor = input.IdAuditor
		} else {
			http.Error(w, helpers.Error("Not localadmin"), http.StatusForbidden)
			return
		}
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

	audit := r.Context().Value("audit").(*datastore.Audit)

	err := h.Datastore.DeleteAuditById(audit.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getAuditSubgroups(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	filter := helpers.GetQueryArgs([][]string{
		[]string{"idAudit", "id_audit"},
		[]string{"idSubgroup", "id_subgroup"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	subgroups, err := h.Datastore.GetAuditSubgroupsByIdAudit(audit.Id, filter)
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

	audit := r.Context().Value("audit").(*datastore.Audit)

	var input struct {
		Subgroups []int64 `json:"subgroups"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		if err := r.ParseForm(); err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
		idsSubgroupsStr := r.PostForm["idSubgroup"]
		if nSubgroups := len(idsSubgroupsStr); nSubgroups != 0 {
			input.Subgroups = make([]int64, nSubgroups)
			for index, idSubgroupStr := range idsSubgroupsStr {
				var idSubgroup int64
				idSubgroup, err := helpers.ParseInt64(idSubgroupStr)
				if err != nil {
					http.Error(w, helpers.Error(err.Error()), 400)
					return
				}
				input.Subgroups[index] = idSubgroup
			}
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

	if len(input.Subgroups) == 0 {
		http.Error(w, helpers.Error("at least one subgroup"), 400)
		return
	}

	err := h.Datastore.SaveAuditSubgroup(audit.Id, audit.IdTemplate, input.Subgroups)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteAuditSubgroups(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	err := h.Datastore.DeleteAuditSubgroupsByIdAudit(audit.Id)
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

func (h *Handler) deleteAuditCriteria(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	err := h.Datastore.DeleteAuditCriterionByIdAudit(audit.Id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getAuditCriterion(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

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

	idCriterion := r.Context().Value("idCriterion").(int64)

	var input struct {
		Value int64 `json:"value"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.Value, err = helpers.ParseInt64(r.PostFormValue("value"))
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

	auditCriterion := datastore.NewAuditCriterion(false)
	auditCriterion.IdAudit = audit.Id
	auditCriterion.IdCriterion = idCriterion
	auditCriterion.Value = zero.IntFrom(input.Value)

	err := h.Datastore.SaveAuditCriterion(auditCriterion)
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

	idCriterion := r.Context().Value("idCriterion").(int64)

	auditCriterion, err := h.Datastore.GetAuditCriterionByIds(audit.Id, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	var input struct {
		Value int64 `json:"value"`
	}
	input.Value = -1

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.Value, err = helpers.ParseInt64(r.PostFormValue("value"))
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

	if input.Value == -1 {
		http.Error(w, helpers.Error("At least one of value, observation"), 400)
		return
	}

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

func (h *Handler) deleteAuditCriterion(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	err := h.Datastore.DeleteAuditCriterionByIds(audit.Id, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getCriterionRemarks(w http.ResponseWriter, r *http.Request) {

	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	auditsCriteriaRemarks, err := h.Datastore.GetRemarksByIdsAuditCriterion(audit.Id, idCriterion)
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

	idCriterion := r.Context().Value("idCriterion").(int64)

	var input struct {
		Observation   string `json:"observation"`
		imageMimetype string `json:"-"`
		imageBytes    []byte `json:"-"`
		imageHash     string `json:"-"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		input.Observation = r.PostFormValue("Observation")
		var err error
		input.imageBytes, input.imageMimetype, input.imageHash, err = helpers.ReadImage(r, "image", helpers.MaxImageFileSize)
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
	remark.ImageMimetype = zero.StringFrom(input.imageMimetype)

	err := h.Datastore.InsertRemark(remark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"id":"` + helpers.Int64ToString(remark.Id) + "}"))
}

func (h *Handler) deleteCriterionRemarks(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	err := h.Datastore.DeleteRemarkByIdsAuditCriterion(audit.Id, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getCriterionRemark(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	idRemark := r.Context().Value("idRemark").(int64)

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

func (h *Handler) updateCriterionRemark(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	idRemark := r.Context().Value("idRemark").(int64)

	var input struct {
		Observation   string `json:"observation"`
		imageMimetype string `json:"-"`
		imageBytes    []byte `json:"-"`
		imageHash     string `json:"-"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		input.Observation = r.PostFormValue("observation")
		var err error
		input.imageBytes, input.imageMimetype, input.imageHash, err = helpers.ReadImage(r, "image", helpers.MaxImageFileSize)
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

	remark, err := h.Datastore.GetRemarkByIdsAuditCriterionRemark(audit.Id, idCriterion, idRemark)
	if remark != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	if input.Observation != "" {
		remark.Observation = zero.StringFrom(input.Observation)
	}

	if input.imageBytes != nil {
		remark.Image = input.imageBytes
	}

	err = h.Datastore.UpdateRemark(remark)
	if remark != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteCriterionRemark(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	idRemark := r.Context().Value("idRemark").(int64)

	err := h.Datastore.DeleteRemarkByIdsAuditCriterionRemark(audit.Id, idCriterion, idRemark)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

func (h *Handler) getCriterionRemarkImage(w http.ResponseWriter, r *http.Request) {
	audit := r.Context().Value("audit").(*datastore.Audit)

	idCriterion := r.Context().Value("idCriterion").(int64)

	idRemark := r.Context().Value("idRemark").(int64)

	urlHash := chi.URLParam(r, "hash")
	if !helpers.ImageHashSizeIsValid(len(urlHash)) {
		http.Error(w, helpers.Error("Invalid image hash"), 400)
		return
	}

	auditCriterionRemark, err := h.Datastore.GetRemarkByIdsAuditCriterionRemark(audit.Id, idCriterion, idRemark)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), http.StatusBadRequest)
		return
	}

	if auditCriterionRemark.ImageMimetype.Valid && auditCriterionRemark.Image != nil && auditCriterionRemark.ImageHash.Valid {
		scheme := "https://"
		if r.TLS == nil {
			scheme = "http://"
		}
		redirectUrl := helpers.FastConcat(scheme, r.Host, "/audits/", helpers.Int64ToString(audit.Id), "/criteria/", helpers.Int64ToString(idCriterion), "/remarks/", helpers.Int64ToString(idRemark), "/image/", auditCriterionRemark.ImageHash.String)
		if err := helpers.ImageCashingControlWriter(w, r, urlHash, auditCriterionRemark.ImageHash.String, auditCriterionRemark.ImageMimetype.String, auditCriterionRemark.Image, redirectUrl); err != nil {
			http.Error(w, helpers.Error(err.Error()), http.StatusInternalServerError)
			return
		}
		return
	} else {
		http.Error(w, helpers.Error("no image"), http.StatusNotFound)
		return
	}
}

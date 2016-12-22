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
	"fmt"
	"os"
	"io"
	"mime/multipart"
	"path/filepath"
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

	//POST /audits/:ida/criteria/:idc/remarks envia {...}  retorna id do remark;
	//GET /audits/:ida/criteria/:idc/remarks/:idr ;
	//GET   /audits/:ida/criteria/:idc/remarks

	router.Post("/:ida/criteria/:idc/remarks", h.createCriterionRemark)
	router.Get("/:ida/criteria/:idc/remarks/:idr", h.getCriterionRemark)
	router.Get("/:ida/criteria/:idc/remarks", helpers.ReplyJson(h.getCriterionRemarks))
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
	_, err = h.Datastore.GetAuditorById(input.IdAuditor)
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

	filter := helpers.GetQueryArgs([][]string{
		[]string{"idAudit", "id_audit"},
		[]string{"idCriterion", "id_criterion"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	auditCriteria, err := h.Datastore.GetAuditCriteria(idAudit, filter)
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

	auditCriterion, err := h.Datastore.GetAuditCriterionById(idAudit, idCriterion)
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

	auditCriterion, err := h.Datastore.GetAuditCriterionById(idAudit, idCriterion)
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

	err = h.Datastore.DeleteAuditCriterionById(idAudit, idCriterion)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}
func (h *Handler) createCriterionRemark(w http.ResponseWriter, r *http.Request) {
//	router.Post("/:ida/criteria/:idc/remarks", h.createCriterionRemark)
	//*******************************erro a receber os paramentros
	idAuditStr := chi.URLParam(r, "ida")
	idCriterionStr := chi.URLParam(r, "idc")
	idAuditStr = "1"
	fmt.Printf("%v \n",idAuditStr)
	fmt.Printf("%v \n",idCriterionStr)
	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idCriterion, err := strconv.ParseInt(idCriterionStr, 10, 64)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//http://stackoverflow.com/questions/28940005/golang-get-multipart-form-data
	//http://stackoverflow.com/questions/25225723/passing-a-image-from-a-html-form-to-go
	err = r.ParseMultipartForm(24 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//////****************falta receber o path
	in, header, err := r.FormFile("file")
	if err!= nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//you probably want to make sure header.Filename is unique and
	// use filepath.Join to put it somewhere else.
	out, err := os.OpenFile(header.Filename, os.O_RDWR | os.O_EXCL , 0644)
	fmt.Printf("abrir ficheiro \n")
	fmt.Printf("ficheiro %v \n",out)
	fmt.Printf("nome ficheiro %v \n",header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, in)
	fmt.Fprintf(w, "Readbytes %v", in)
	//http://stackoverflow.com/questions/22945486/golang-converting-image-image-to-byte
	//leitura do ficheiro
	var buf []byte
	_,err = out.Read(buf)
	fmt.Printf("convert a imagem para bytes \n")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Readbytes %v, Handler %v", buf, header.Filename)
	fmt.Printf("inicio do remarks \n")
	//observation
	observation := r.FormValue("observation")
	remark := datastore.NewRemark(false)
	remark.Image = buf
	remark.IdAudit = idAudit
	remark.IdCriterion = idCriterion
	remark.Observation = zero.StringFrom(observation)
	h.Datastore.InsertRemark(remark)
}

func (h *Handler) getCriterionRemark(w http.ResponseWriter, r *http.Request){
//router.Get("/:ida/criteria/:idc/remarks/:idr", h.getCriterionRemark)

	idAuditStr:= chi.URLParam(r, "ida")
	idCriteriaStr := chi.URLParam(r, "idc")
	idRemarkStr := chi.URLParam(r, "idr")

	fmt.Printf("%v \n",idAuditStr)
	fmt.Printf("%v \n",idCriteriaStr)
	fmt.Printf("%v \n",idRemarkStr)

	idAudit, err := strconv.ParseInt(idAuditStr, 10, 64)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idCriterion, err := strconv.ParseInt(idCriteriaStr, 10, 64)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idRemark, err := strconv.ParseInt(idRemarkStr, 10, 64)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	auditsCriteriaRemarks, err := h.Datastore.GetRemarkByAuditCriterionIds( idAudit, idCriterion, idRemark)
//convert para json
/*	auditsCriteriaRemarksSlice, err := json.Marshal(auditsCriteriaRemarks)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	_, err = w.Write(auditsCriteriaRemarksSlice)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}*/
	//****** https://github.com/gebi/go-fileupload-example/blob/master/main.go
	//****** http://stackoverflow.com/questions/22945486/golang-converting-image-image-to-byte


	//******

	//para multiform data
	//https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
	out, err := os.OpenFile("tempFile", os.O_RDWR | os.O_EXCL , 0644)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	fmt.Fprintf(w, "Readbytes %v", out)
	//http://stackoverflow.com/questions/22945486/golang-converting-image-image-to-byte
	//leitura do ficheiro
	_,err = out.Write(auditsCriteriaRemarks.Image)
	writerMultiPart := multipart.NewWriter(w)
	part, err := writerMultiPart.CreateFormFile("tempFile", filepath.Base("/tempFiles"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(part, out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = writerMultiPart.WriteField("observation", auditsCriteriaRemarks.Observation.String)

	err = writerMultiPart.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//convert para byes
	auditsCriteriaRemarksSlice, err := json.Marshal(writerMultiPart)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	//responde
	_, err = w.Write(auditsCriteriaRemarksSlice)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

}
func (h *Handler) getCriterionRemarks(w http.ResponseWriter, r *http.Request){


	idAudits := chi.URLParam(r, "ida")
	idCriteria := chi.URLParam(r, "idc")

	idaudit, err := strconv.ParseInt(idAudits, 10, 64)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	idcriterion, err := strconv.ParseInt(idCriteria, 10, 64)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	auditsCriteriaRemarks, err := h.Datastore.GetRemarksByAuditCriterionIds(idaudit, idcriterion)


	auditsCriteriaRemarksSlice, err := json.Marshal(auditsCriteriaRemarks)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	_, err = w.Write(auditsCriteriaRemarksSlice)
	if err!=nil{
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	//****** https://github.com/gebi/go-fileupload-example/blob/master/main.go
	//****** http://stackoverflow.com/questions/22945486/golang-converting-image-image-to-byte


	//******



}

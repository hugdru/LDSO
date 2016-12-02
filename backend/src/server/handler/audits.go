package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"
	"time"
	"fmt"
	"gopkg.in/guregu/null.v3/zero"
)

func (h *Handler) auditsRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getAudits))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createAudit)))
	router.Get("/:id", helpers.ReplyJson(h.getAudit))
	router.Put("/:id", helpers.RequestJson(helpers.ReplyJson(h.updateAudit)))
	router.Delete("/:id", helpers.ReplyJson(h.deleteAudit))
}

func (h *Handler) getAudits(w http.ResponseWriter, r *http.Request) {
	var limit, offset int = 10, 0
	var err error

	limitString := r.FormValue("limit")
	offsetString := r.FormValue("offset")

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
	}
	if offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
	}

	if limit <= 0 || offset < 0 || limit > 100 {
		http.Error(w, helpers.Error("0<limit<=100 && offset > 0"), 400)
		return
	}

	audits, err := h.Datastore.GetAudits(limit,offset);
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
	fmt.Println(input.IdAuditor)
	fmt.Println(input.IdTemplate)
	fmt.Println(input.IdProperty)
	if input.IdProperty == 0 || input.IdAuditor == 0 || input.IdTemplate == 0 {
		http.Error(w, helpers.Error("The audits must have idProperty, idAuditor and idTemplate"), 400)
		return
	}

	_, err = h.Datastore.GetTemplateById(input.IdProperty)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_, err = h.Datastore.GetTemplateById(input.IdAuditor)
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
		IdAuditor  int64
		Rating       int64
		Observation  string
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
		audit.Observation =  zero.StringFrom(input.Observation)
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
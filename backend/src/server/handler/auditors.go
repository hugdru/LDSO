package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"strconv"

	"fmt"
	"time"
)

func (h *Handler) auditorsRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getAuditors))
	router.Post("/", helpers.RequestJson(helpers.ReplyJson(h.createAuditor)))
	router.Get("/:id", helpers.ReplyJson(h.getAuditor))
	router.Put("/:id", helpers.RequestJson(helpers.ReplyJson(h.updateAuditor)))
	router.Delete("/:id", helpers.ReplyJson(h.deleteAuditor))
}

///******begin getAuditors************/////
func (h *Handler) getAuditors(w http.ResponseWriter, r *http.Request) {
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

	auditors, err := h.Datastore.GetAuditors(limit, offset)
	fmt.Println(auditors)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	auditorsSlice, err := json.Marshal(auditors)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditorsSlice)
}

///******end getAuditors************/////

///******begin createAuditor************/////
func (h *Handler) createAuditor(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	if decoder == nil {
		http.Error(w, helpers.Error("JSON decoder failed"), 500)
		return
	}

	var input struct {
		IdCountry int64
		Name      string
		Email     string
		Username  string
		Password  string
	}
	input.IdCountry = 0
	input.Name = ""
	input.Email = ""
	input.Username = ""
	input.Password = ""

	err := decoder.Decode(&input)

	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	if input.IdCountry == 0 {
		http.Error(w, helpers.Error("The auditor must have idCountry"), 400)
		return
	}
	_, err = h.Datastore.GetCountryById(input.IdCountry)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	entity := datastore.NewEntity(false)
	entity.IdCountry = input.IdCountry
	entity.Name = input.Name
	entity.Email = input.Email
	entity.Username = input.Username
	entity.Password = input.Password
	entity.CreatedDate = time.Now().UTC()
	auditor := datastore.NewAuditor(false)
	auditor.IdEntity = entity.Id
	err = h.Datastore.SaveEntity(entity)

	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	//grava o entity
	entitySlice, err := json.Marshal(entity)

	err = h.Datastore.SaveAuditor(auditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(entitySlice)
	//grava auditor
	auditorSlice, err := json.Marshal(auditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditorSlice)

}

///******end createAuditor************/////

///******begin getAuditor************/////

func (h *Handler) getAuditor(w http.ResponseWriter, r *http.Request) {
	idAuditor := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAuditor, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	auditor, err := h.Datastore.GetAuditorById(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	auditorSlice, err := json.Marshal(auditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(auditorSlice)
}

///******end getAuditor************/////

///******begin updateAuditor************/////

func (h *Handler) updateAuditor(w http.ResponseWriter, r *http.Request) {
	/*idAuditor := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idAuditor, 10, 64)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		//auditor, err := h.Datastore.GetAuditorById(id)
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
	/*
		if input.IdAuditor != 0 {
			auditor.IdAuditor = input.IdAuditor
		}


		if input.Observation != "" {
			audit.Observation =  zero.StringFrom(input.Observation)
		}

		err = h.Datastore.SaveAudit(auditor)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 500)
			return
		}
	*/
}

///******end updateAuditor************/////

///******begin deleteAuditor************/////
func (h *Handler) deleteAuditor(w http.ResponseWriter, r *http.Request) {
	idAuditor := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAuditor, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	auditor, err := h.Datastore.GetAuditorById(id)

	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	err = h.Datastore.DeleteAuditor(auditor)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

}

///******end deleteAuditor************/////

///******end Page************/////
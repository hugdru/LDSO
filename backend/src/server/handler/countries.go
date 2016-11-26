package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) countriesRoutes(router chi.Router) {
	router.Get("/", helpers.ReplyJson(h.getCountries))
	router.Get("/:id", helpers.ReplyJson(h.getCountry))
}

func (h *Handler) getCountry(w http.ResponseWriter, r *http.Request) {
	idCountry := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idCountry, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error("Bad Country id"), 400)
		return
	}
	country, err := h.Datastore.GetCountryById(id)
	if err != nil {
		http.Error(w, helpers.Error("Invalid Country id"), 400)
		return
	}
	countrySlice, err := json.Marshal(country)
	if err != nil {
		http.Error(w, helpers.Error("Failed converting Country to JSON"), 500)
		return
	}
	w.Write(countrySlice)
}

func (h *Handler) getCountries(w http.ResponseWriter, r *http.Request) {
	country, err := h.Datastore.GetCountries()
	if err != nil {
		http.Error(w, helpers.Error("Failed grabbing Countries"), 500)
		return
	}
	countrySlice, err := json.Marshal(country)
	if err != nil {
		http.Error(w, helpers.Error("Failed converting Country to JSON"), 500)
		return
	}
	w.Write(countrySlice)
}

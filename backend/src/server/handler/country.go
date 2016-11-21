package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"strconv"
)

const countryMain = "/country"

func (h *Handler) CountryMain() string {
	return countryMain
}

func (h *Handler) CountrySubroutes(router chi.Router) {
	router.Get("/:id", helpers.ReplyJson(h.getCountry))
}

func (h *Handler) getCountry(w http.ResponseWriter, r *http.Request) {
	idQuery := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idQuery, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error("Bad id"), 400)
		return
	}
	country, err := h.Datastore.GetCountryById(id)
	if err != nil {
		http.Error(w, "Invalid country id", 400)
		return
	}
	countrySlice, err := json.Marshal(country)
	if err != nil {
		http.Error(w, helpers.Error("Failed converting country to JSON"), 500)
		return
	}
	w.Write(countrySlice)
}

package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) addressesRoutes(router chi.Router) {
	router.Get("/:id", helpers.ReplyJson(h.getAddress))
}

func (h *Handler) getAddress(w http.ResponseWriter, r *http.Request) {
	idAddress := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAddress, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	address, err := h.Datastore.GetAddressByIdWithForeign(id)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	addressSlice, err := json.Marshal(address)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(addressSlice)
}

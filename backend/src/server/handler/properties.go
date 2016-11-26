package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) propertiesRoutes(router chi.Router) {
	router.Get("/:id", helpers.ReplyJson(h.GetProperty))
}

func (h *Handler) GetProperty(w http.ResponseWriter, r *http.Request) {
	idProperty := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idProperty, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error("Bad Property id"), 400)
		return
	}
	property, err := h.Datastore.GetPropertyById(id)
	if err != nil {
		http.Error(w, helpers.Error("Invalid Property id"), 400)
		return
	}
	propertySlice, err := json.Marshal(property)
	if err != nil {
		http.Error(w, helpers.Error("Failed converting Property to JSON"), 500)
		return
	}
	w.Write(propertySlice)
}

package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"strconv"
)

func (h *Handler) clientsRoutes(router chi.Router) {
	router.Get("/:id", helpers.ReplyJson(h.getClient))
}

func (h *Handler) getClient(w http.ResponseWriter, r *http.Request) {
	idClient := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idClient, 10, 64)
	if err != nil {
		http.Error(w, helpers.Error("Bad Client id"), 400)
		return
	}
	client, err := h.Datastore.GetClientById(id)
	if err != nil {
		http.Error(w, helpers.Error("Invalid Client id"), 400)
		return
	}
	clientSlice, err := json.Marshal(client)
	if err != nil {
		http.Error(w, helpers.Error("Failed converting Client to JSON"), 500)
		return
	}
	w.Write(clientSlice)
}

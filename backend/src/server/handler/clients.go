package handler

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"server/handler/sessionData"
)

func (h *Handler) clientsRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getClients))
	router.Post("/", decorators.OnlySuperadminsOrLocaladmins(decorators.ReplyJson(h.createClient)))
	router.Route("/:idc", h.clientRoutes)
}

func (h *Handler) clientRoutes(router chi.Router) {
	router.Use(h.clientContext)
	router.Get("/", decorators.ReplyJson(h.getClient))
	router.Get("/image/:hash", h.getClientImage)
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrClients(decorators.ReplyJson(h.updateClient)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrClients(decorators.ReplyJson(h.deleteClient)))
}

func (h *Handler) clientContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idClientStr := chi.URLParam(r, "idc")
		idClient, err := helpers.ParseInt64(idClientStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		esd, err := sessionData.GetSessionData(r)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		restricted := true
		switch esd.Role {
		case sessionData.Superadmin, sessionData.Localadmin, sessionData.Auditor:
			restricted = false
		case sessionData.Client:
			if esd.Id == idClient {
				restricted = false
			}
		}

		client, err := h.Datastore.GetClientById(idClient, true, restricted)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		if r.Method != http.MethodGet {
			if esd.Id != client.IdEntity && esd.Role != sessionData.Superadmin && esd.Role != sessionData.Localadmin {
				http.Error(w, helpers.Error("Not the owner of the account"), http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), "client", client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getClients(w http.ResponseWriter, r *http.Request) {

	esd, err := sessionData.GetSessionData(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	restricted := true
	switch esd.Role {
	case sessionData.Superadmin, sessionData.Localadmin, sessionData.Auditor:
		restricted = false
	}

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	clients, err := h.Datastore.GetClients(limit, offset, true, restricted)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	clientsSlice, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(clientsSlice)
}

func (h *Handler) createClient(w http.ResponseWriter, r *http.Request) {

	var input struct {
		IdEntity int64 `json:"IdEntity"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		var err error
		input.IdEntity, err = helpers.ParseInt64(r.PostFormValue("IdEntity"))
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

	newClient := datastore.NewClient(false)
	err := newClient.MustSet(input.IdEntity)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveClient(newClient)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	clientSlice, err := json.Marshal(newClient)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(clientSlice)
}

func (h *Handler) getClient(w http.ResponseWriter, r *http.Request) {

	client := r.Context().Value("client").(*datastore.Client)

	clientSlice, err := json.Marshal(client)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(clientSlice)
}

func (h *Handler) getClientImage(w http.ResponseWriter, r *http.Request) {

	client := r.Context().Value("client").(*datastore.Client)
	entity := client.Entity

	urlHash := chi.URLParam(r, "hash")

	scheme := "https://"
	if r.TLS == nil {
		scheme = "http://"
	}
	redirectUrl := helpers.FastConcat(scheme, r.Host, "/clients/", helpers.Int64ToString(entity.Id), "/image/", entity.ImageHash.String)
	getEntityImage(w, r, entity, urlHash, redirectUrl)
}

func (h *Handler) updateClient(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value("client").(*datastore.Client)
	if !h.updateEntity(w, r, nil, client.Entity) {
		return
	}
}

func (h *Handler) deleteClient(w http.ResponseWriter, r *http.Request) {

	client := r.Context().Value("client").(*datastore.Client)

	err := h.Datastore.DeleteClient(client)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

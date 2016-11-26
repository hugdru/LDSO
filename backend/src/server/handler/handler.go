package handler

import (
	"github.com/pressly/chi"
	"server/datastore"
)

type Handler struct {
	Datastore *datastore.Datastore
}

type routing struct {
	entrypoint string
	fn         func(chi.Router)
}

func (h *Handler) Init(mux *chi.Mux) {
	mux.Route("/countries", h.countriesRoutes)
	mux.Route("/properties", h.propertiesRoutes)
	mux.Route("/addresses", h.addressesRoutes)
	mux.Route("/clients", h.clientsRoutes)
	mux.Route("/templates", h.templatesRoutes)
}

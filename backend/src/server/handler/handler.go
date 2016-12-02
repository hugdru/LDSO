package handler

import (
	"github.com/pressly/chi"
	"server/datastore"
)

type Handler struct {
	Datastore *datastore.Datastore
}

func (h *Handler) Init(mux *chi.Mux) {
	mux.Route("/countries", h.countriesRoutes)
	mux.Route("/properties", h.propertiesRoutes)
	mux.Route("/addresses", h.addressesRoutes)
	mux.Route("/clients", h.clientsRoutes)
	mux.Route("/templates", h.templatesRoutes)
	mux.Route("/maingroups", h.maingroupsRoutes)
	mux.Route("/subgroups", h.subgroupsRoutes)
	mux.Route("/legislations", h.legislationsRoutes)
	mux.Route("/criteria", h.criteriaRoutes)
	mux.Route("/audits", h.auditsRoutes)
}

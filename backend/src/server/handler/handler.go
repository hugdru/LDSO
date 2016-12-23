package handler

import (
	"github.com/pressly/chi"
	"server/datastore"
	"server/handler/sessionData"
)

type Handler struct {
	Datastore *datastore.Datastore
}

func (h *Handler) Init(mux *chi.Mux) {
	sessionData.GobRegister()
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
	mux.Route("/accessibilities", h.accessibilitiesRoutes)
	mux.Route("/auditors", h.auditorsRoutes)
	mux.Route("/entities", h.entitiesRoutes)
}

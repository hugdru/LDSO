package handler

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
	"context"
	"gopkg.in/guregu/null.v3/zero"
	"server/handler/sessionData"
	"server/datastore"
)

func (h *Handler) propertiesRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getProperties))
	router.Post("/", decorators.OnlyClients(h.createProperty))
	router.Route("/:idp", h.propertyRoutes)
}

func (h *Handler) propertyRoutes(router chi.Router) {
	router.Use(h.propertyContext)
	router.Get("/", decorators.ReplyJson(h.getProperty))
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrClients(h.updateProperty))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrClients(h.deleteProperty))
	router.Get("/address", decorators.ReplyJson(h.getAddress))
	router.Put("/address", decorators.OnlySuperadminsOrLocaladminsOrClients(decorators.ReplyJson(h.updateAddress)))
}

func (h *Handler) propertyContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idPropertyStr := chi.URLParam(r, "idp")
		idProperty, err := helpers.ParseInt64(idPropertyStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		property, err := h.Datastore.GetPropertyByIdWithForeign(idProperty)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		if r.Method != http.MethodGet {
			entitySessionData, err := sessionData.GetSessionData(r)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 400)
				return
			}
			owners, err := h.Datastore.GetPropertyClientsByIdProperty(idProperty)
			if err != nil {
				http.Error(w, helpers.Error(err.Error()), 400)
				return
			}
			found := false
			for _, owner := range owners {
				if owner.IdEntity == entitySessionData.Id {
					found = true
					break
				}
			}
			if !found && entitySessionData.Role != sessionData.Superadmin && entitySessionData.Role != sessionData.Localadmin {
				http.Error(w, helpers.Error("Not the owner of the property"), http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), "property", property)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getProperties(w http.ResponseWriter, r *http.Request) {

	limit, offset, err := helpers.PaginationParse(r)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	properties, err := h.Datastore.GetProperties(limit, offset, filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	propertiesSlice, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(propertiesSlice)
}

func (h *Handler) createProperty(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getProperty(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)

	propertySlice, err := json.Marshal(property)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(propertySlice)
}

func (h *Handler) updateProperty(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)

}

func (h *Handler) deleteProperty(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)
}

func (h *Handler) getAddress(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)

	address, err := h.Datastore.GetAddressByIdWithForeign(property.IdAddress)
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

func (h *Handler) updateAddress(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)

	var input struct {
		IdCountry    int64  `json:"idCountry"`
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		AddressLine3 string `json:"addressLine3"`
		TownCity     string `json:"townCity"`
		County       string `json:"county"`
		Postcode     string `json:"postcode"`
		Latitude     string `json:"latitude"`
		Longitude    string `json:"longitude"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		var err error
		input.IdCountry, err = helpers.ParseInt64(r.PostFormValue("idCountry"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.AddressLine1 = r.PostFormValue("addressLine1")
		input.AddressLine2 = r.PostFormValue("addressLine2")
		input.AddressLine3 = r.PostFormValue("addressLine3")
		input.TownCity = r.PostFormValue("townCity")
		input.County = r.PostFormValue("county")
		input.Postcode = r.PostFormValue("postcode")
		input.Latitude = r.PostFormValue("latitude")
		input.Longitude = r.PostFormValue("longitude")
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

	if input.IdCountry == 0 && input.AddressLine1 == "" && input.AddressLine2 == "" && input.AddressLine3 == "" &&
		input.TownCity == "" && input.County == "" && input.Postcode == "" && input.Latitude == "" && input.Longitude == "" {
		http.Error(w, helpers.Error("Change at least one address value"), 400)
		return
	}

	address, err := h.Datastore.GetAddressByIdWithForeign(property.IdAddress)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if input.IdCountry != 0 {
		address.IdCountry = input.IdCountry
	}

	if input.AddressLine1 != "" {
		address.AddressLine1 = input.AddressLine1
	}

	if input.AddressLine2 != "" {
		address.AddressLine2 = zero.StringFrom(input.AddressLine2)
	}

	if input.AddressLine3 != "" {
		address.AddressLine3 = zero.StringFrom(input.AddressLine3)
	}

	if input.TownCity != "" {
		address.TownCity = zero.StringFrom(input.TownCity)
	}

	if input.County != "" {
		address.County = zero.StringFrom(input.County)
	}

	if input.Postcode != "" {
		address.Postcode = zero.StringFrom(input.Postcode)
	}

	if input.Latitude != "" {
		address.Latitude = zero.StringFrom(input.Latitude)
	}

	if input.Longitude != "" {
		address.Longitude = zero.StringFrom(input.Longitude)
	}

	if h.Datastore.SaveAddress(address) != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

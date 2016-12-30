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

func (h *Handler) propertiesRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getProperties))
	router.Post("/", decorators.OnlyClients(decorators.ReplyJson(h.createProperty)))
	router.Route("/:idp", h.propertyRoutes)
}

func (h *Handler) propertyRoutes(router chi.Router) {
	router.Use(h.propertyContext)
	router.Get("/", decorators.ReplyJson(h.getProperty))
	router.Put("/", decorators.OnlySuperadminsOrLocaladminsOrClients(decorators.ReplyJson(h.updateProperty)))
	router.Delete("/", decorators.OnlySuperadminsOrLocaladminsOrClients(decorators.ReplyJson(h.deleteProperty)))
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

		esd, err := sessionData.GetSessionData(r)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		owners, err := h.Datastore.GetPropertyClientsByIdProperty(idProperty, false, false)
		isOwner := false
		for _, owner := range owners {
			if owner.IdEntity == esd.Id {
				isOwner = true
				break
			}
		}

		restricted := true
		switch esd.Role {
		case sessionData.Superadmin, sessionData.Localadmin, sessionData.Auditor:
			restricted = false
		case sessionData.Client:
			if isOwner {
				restricted = false
			}
		}

		if r.Method != http.MethodGet {
			if !isOwner && esd.Role != sessionData.Superadmin && esd.Role != sessionData.Localadmin {
				http.Error(w, helpers.Error("Not the owner of the property"), http.StatusForbidden)
				return
			}
		}

		property, err := h.Datastore.GetPropertyByIdWithAddressTagsOwners(idProperty, true, restricted)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "property", property)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getProperties(w http.ResponseWriter, r *http.Request) {

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

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	properties, err := h.Datastore.GetPropertiesWithAddressTagsOwners(limit, offset, filter, true, restricted)
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

	var input struct {
		Name    string `json:"name"`
		Details string `json:"details"`
		Address struct {
			IdCountry    int64  `json:"idCountry"`
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			AddressLine3 string `json:"addressLine3"`
			TownCity     string `json:"townCity"`
			County       string `json:"county"`
			Postcode     string `json:"postcode"`
			Latitude     string `json:"latitude"`
			Longitude    string `json:"longitude"`
		} `json:"address"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		input.Name = r.PostFormValue("name")
		input.Details = r.PostFormValue("details")
		var err error
		input.Address.IdCountry, err = helpers.ParseInt64(r.PostFormValue("address.idCountry"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.Address.AddressLine1 = r.PostFormValue("address.addressLine1")
		input.Address.AddressLine2 = r.PostFormValue("address.addressLine2")
		input.Address.AddressLine3 = r.PostFormValue("address.addressLine3")
		input.Address.TownCity = r.PostFormValue("address.townCity")
		input.Address.County = r.PostFormValue("address.county")
		input.Address.Postcode = r.PostFormValue("address.postcode")
		input.Address.Latitude = r.PostFormValue("address.latitude")
		input.Address.Longitude = r.PostFormValue("address.longitude")
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

	newProperty := datastore.NewProperty(false)
	err := newProperty.MustSet(input.Name, input.Details)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	newProperty.Address = datastore.NewAddress(false)
	err = newProperty.Address.MustSet(input.Address.IdCountry, input.Address.AddressLine1)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = newProperty.Address.OptionalSetIfNotEmptyOrNil(
		input.Address.AddressLine2,
		input.Address.AddressLine3,
		input.Address.TownCity,
		input.Address.County,
		input.Address.Postcode,
		input.Address.Latitude,
		input.Address.Longitude,
	)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SavePropertyWithAddress(newProperty)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	propertySlice, err := json.Marshal(newProperty)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(propertySlice)
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

	var input struct {
		Name    string `json:"name"`
		Details string `json:"details"`
		Address struct {
			IdCountry    int64  `json:"idCountry"`
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			AddressLine3 string `json:"addressLine3"`
			TownCity     string `json:"townCity"`
			County       string `json:"county"`
			Postcode     string `json:"postcode"`
			Latitude     string `json:"latitude"`
			Longitude    string `json:"longitude"`
		} `json:"address"`
	}

	contentType := helpers.GetContentType(r.Header.Get("Content-type"))
	switch contentType {
	case "multipart/form-data":
		input.Name = r.PostFormValue("name")
		input.Details = r.PostFormValue("details")
		var err error
		input.Address.IdCountry, err = helpers.ParseInt64(r.PostFormValue("address.idCountry"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.Address.AddressLine1 = r.PostFormValue("address.addressLine1")
		input.Address.AddressLine2 = r.PostFormValue("address.addressLine2")
		input.Address.AddressLine3 = r.PostFormValue("address.addressLine3")
		input.Address.TownCity = r.PostFormValue("address.townCity")
		input.Address.County = r.PostFormValue("address.county")
		input.Address.Postcode = r.PostFormValue("address.postcode")
		input.Address.Latitude = r.PostFormValue("address.latitude")
		input.Address.Longitude = r.PostFormValue("address.longitude")
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

	if input.Name == "" && input.Details == "" && input.Address.IdCountry == 0 && input.Address.AddressLine1 == "" &&
		input.Address.AddressLine2 == "" && input.Address.AddressLine3 == "" && input.Address.TownCity == "" &&
		input.Address.County == "" && input.Address.Postcode == "" && input.Address.Latitude == "" && input.Address.Longitude == "" {
		http.Error(w, helpers.Error("At least one argument"), 400)
		return
	}

	err := property.AllSetIfNotEmptyOrNil(input.Name, input.Details)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
	err = property.Address.AllSetIfNotEmptyOrNil(
		input.Address.IdCountry,
		input.Address.AddressLine1,
		input.Address.AddressLine2,
		input.Address.AddressLine3,
		input.Address.TownCity,
		input.Address.County,
		input.Address.Postcode,
		input.Address.Latitude,
		input.Address.Longitude,
	)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SavePropertyWithAddress(property)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	propertySlice, err := json.Marshal(property)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(propertySlice)

}

func (h *Handler) deleteProperty(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)
	err := h.Datastore.DeletePropertyWithAddress(property)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) getAddress(w http.ResponseWriter, r *http.Request) {

	property := r.Context().Value("property").(*datastore.Property)

	address, err := h.Datastore.GetAddressByIdWithCountry(property.IdAddress)
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
		input.IdCountry, err = helpers.ParseInt64(r.PostFormValue("address.idCountry"))
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		input.AddressLine1 = r.PostFormValue("address.addressLine1")
		input.AddressLine2 = r.PostFormValue("address.addressLine2")
		input.AddressLine3 = r.PostFormValue("address.addressLine3")
		input.TownCity = r.PostFormValue("address.townCity")
		input.County = r.PostFormValue("address.county")
		input.Postcode = r.PostFormValue("address.postcode")
		input.Latitude = r.PostFormValue("address.latitude")
		input.Longitude = r.PostFormValue("address.longitude")
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

	err := property.Address.AllSetIfNotEmptyOrNil(
		input.IdCountry,
		input.AddressLine1,
		input.AddressLine2,
		input.AddressLine3,
		input.TownCity,
		input.County,
		input.Postcode,
		input.Latitude,
		input.Longitude,
	)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	if h.Datastore.SaveAddress(property.Address) != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}
}

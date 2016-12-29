package handler

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/datastore"
	"server/handler/helpers"
	"server/handler/helpers/decorators"
)

func (h *Handler) accessibilitiesRoutes(router chi.Router) {
	router.Get("/", decorators.ReplyJson(h.getAccessibilities))
	router.Post("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.createAccessibility)))
	router.Route("/:ida", h.accessibilityRoutes)
}

func (h *Handler) accessibilityRoutes(router chi.Router) {
	router.Use(h.accessibilityContext)
	router.Get("/", decorators.ReplyJson(h.getAccessibility))
	router.Put("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.updateAccessibility)))
	router.Delete("/", decorators.OnlySuperadmins(decorators.ReplyJson(h.deleteAccessibility)))
}

func (h *Handler) accessibilityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAccessibilityStr := chi.URLParam(r, "ida")
		idAccessibility, err := helpers.ParseInt64(idAccessibilityStr)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}
		accessibility, err := h.Datastore.GetAccessibilityById(idAccessibility)
		if err != nil {
			http.Error(w, helpers.Error(err.Error()), 400)
			return
		}

		ctx := context.WithValue(r.Context(), "accessibility", accessibility)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getAccessibilities(w http.ResponseWriter, r *http.Request) {

	filter := helpers.GetQueryArgs([][]string{
		[]string{"id"},
		[]string{"name"},
	}, r)
	if filter == nil {
		http.Error(w, helpers.Error("Failed to create filter"), 500)
		return
	}

	accessibilities, err := h.Datastore.GetAccessibilities(filter)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	accessibilitiesSlice, err := json.Marshal(accessibilities)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
	w.Write(accessibilitiesSlice)
}

func (h *Handler) createAccessibility(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string `json:"name"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		input.Name = r.PostFormValue("name")
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

	accessibility := datastore.NewAccessibility(false)
	err := accessibility.MustSet(input.Name)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveAccessibility(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	accessibilitySlice, err := json.Marshal(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(accessibilitySlice)
}

func (h *Handler) getAccessibility(w http.ResponseWriter, r *http.Request) {

	accessibility := r.Context().Value("accessibility").(*datastore.Accessibility)

	accessibilitySlice, err := json.Marshal(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}

	w.Write(accessibilitySlice)
}

func (h *Handler) updateAccessibility(w http.ResponseWriter, r *http.Request) {

	accessibility := r.Context().Value("accessibility").(*datastore.Accessibility)

	var input struct {
		Name string `json:"name"`
	}

	switch helpers.GetContentType(r.Header.Get("Content-type")) {
	case "multipart/form-data":
		input.Name = r.PostFormValue("name")
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

	if input.Name == "" {
		http.Error(w, helpers.Error("name missing"), 400)
		return
	}

	err := accessibility.AllSetIfNotEmptyOrNil(input.Name)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 400)
		return
	}

	err = h.Datastore.SaveAccessibility(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

func (h *Handler) deleteAccessibility(w http.ResponseWriter, r *http.Request) {

	accessibility := r.Context().Value("accessibility").(*datastore.Accessibility)

	err := h.Datastore.DeleteAccessibility(accessibility)
	if err != nil {
		http.Error(w, helpers.Error(err.Error()), 500)
		return
	}
}

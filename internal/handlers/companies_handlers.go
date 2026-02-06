package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

func (h *Handler) companiesHandlerGET(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if id, ok := idFromQuery(r); ok {
		c, err := h.companies.GetByID(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, c, http.StatusOK)
		return
	}
	list, err := h.companies.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, list, http.StatusOK)
}

func (h *Handler) companiesHandlerPOST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var c models.Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Set user_id from authenticated user
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	c.UserID = &claims.UserID
	if err := h.companies.Create(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusCreated)
}

func (h *Handler) companiesHandlerPUT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := idFromQuery(r)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	var c models.Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	// Prevent user from changing the user_id field on update
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	c.UserID = &claims.UserID
	if err := h.companies.Update(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *Handler) companiesHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := idFromQuery(r)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := h.companies.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

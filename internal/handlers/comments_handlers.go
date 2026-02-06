package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

func (h *Handler) commentsHandlerGET(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if id, ok := idFromQuery(r); ok {
		c, err := h.comments.GetByID(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, c, http.StatusOK)
		return
	}
	list, err := h.comments.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, list, http.StatusOK)
}

func (h *Handler) commentsHandlerPOST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.comments.Create(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusCreated)
}

func (h *Handler) commentsHandlerPUT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := idFromQuery(r)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	if err := h.comments.Update(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *Handler) commentsHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := idFromQuery(r)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := h.comments.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

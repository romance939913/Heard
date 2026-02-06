package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

func (h *Handler) postsHandlerGET(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if id, ok := idFromQuery(r); ok {
		p, err := h.posts.GetByID(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, p, http.StatusOK)
		return
	}
	list, err := h.posts.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, list, http.StatusOK)
}

func (h *Handler) postsHandlerPOST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.posts.Create(ctx, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p, http.StatusCreated)
}

func (h *Handler) postsHandlerPUT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := idFromQuery(r)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = id
	if err := h.posts.Update(ctx, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p, http.StatusOK)
}

func (h *Handler) postsHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := idFromQuery(r)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := h.posts.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

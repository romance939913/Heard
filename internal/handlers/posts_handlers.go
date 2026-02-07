package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

func (h *Handler) postsHandlerGET(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if id, ok := idFromQuery(req); ok {
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

func (h *Handler) postsHandlerPOST(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var p models.Post
	
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Set user_id from authenticated user
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserID = claims.UserID

	// Reset upvotes and downvotes to 0 (cannot be set by client)
	p.Upvotes = 0
	p.Downvotes = 0
	if err := h.posts.Create(ctx, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p, http.StatusCreated)
}

func (h *Handler) postsHandlerPUT(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := idFromQuery(req)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// Get the post to verify ownership
	existing, err := h.posts.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}
	// Verify the user owns this post
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if existing.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var p models.Post
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = id
	p.UserID = claims.UserID
	if err := h.posts.Update(ctx, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p, http.StatusOK)
}

func (h *Handler) postsHandlerDELETE(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := idFromQuery(req)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// Get the post to verify ownership
	existing, err := h.posts.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}
	// Verify the user owns this post
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if existing.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.posts.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

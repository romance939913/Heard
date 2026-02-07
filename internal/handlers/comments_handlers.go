package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

func (h *Handler) commentsHandlerGET(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if id, ok := idFromQuery(req); ok {
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

func (h *Handler) commentsHandlerPOST(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var c models.Comment
	if err := json.NewDecoder(req.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Set user_id from authenticated user
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	c.UserID = claims.UserID
	if err := h.comments.Create(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusCreated)
}

func (h *Handler) commentsHandlerPUT(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := idFromQuery(req)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// Get the comment to verify ownership
	existing, err := h.comments.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	}
	// Verify the user owns this comment
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if existing.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var c models.Comment
	if err := json.NewDecoder(req.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	c.UserID = claims.UserID
	if err := h.comments.Update(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusOK)
}

func (h *Handler) commentsHandlerDELETE(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := idFromQuery(req)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// Get the comment to verify ownership
	existing, err := h.comments.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	}
	// Verify the user owns this comment
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if existing.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.comments.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

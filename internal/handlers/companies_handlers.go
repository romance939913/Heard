package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

func (h *Handler) companiesHandlerGET(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if id, ok := idFromQuery(req); ok {
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

func (h *Handler) companiesHandlerPOST(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var c models.Company
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
	c.UserID = &claims.UserID
	if err := h.companies.Create(ctx, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c, http.StatusCreated)
}

func (h *Handler) companiesHandlerPATCH(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := idFromQuery(req)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// Get the existing company to verify ownership and preserve unmodified fields
	existing, err := h.companies.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}
	// Verify the user owns this company
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if existing.UserID == nil || *existing.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// Decode only the fields provided in the request
	var updates map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Apply updates only to provided fields
	if name, ok := updates["name"].(string); ok && name != "" {
		existing.Name = name
	}
	if desc, ok := updates["description"].(string); ok {
		existing.Description = &desc
	}
	if parentID, ok := updates["parent_company_id"].(float64); ok {
		id := int(parentID)
		existing.ParentCompanyID = &id
	}
	if industry, ok := updates["industry"].(string); ok {
		existing.Industry = &industry
	}
	if subIndustry, ok := updates["sub_industry"].(string); ok {
		existing.SubIndustry = &subIndustry
	}
	if hq, ok := updates["headquarters"].(string); ok {
		existing.Headquarters = &hq
	}
	if dateInc, ok := updates["date_incorporated"].(string); ok {
		existing.DateIncorporated = &dateInc
	}
	// user_id cannot be changed
	if err := h.companies.Update(ctx, existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

func (h *Handler) companiesHandlerDELETE(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := idFromQuery(req)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// Get the company to verify ownership
	existing, err := h.companies.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}
	// Verify the user owns this company
	claims, err := GetUserClaimsFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if existing.UserID == nil || *existing.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.companies.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

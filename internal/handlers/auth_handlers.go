package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brennanromance/heard/internal/models"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token   string       `json:"token"`
	User    *models.User `json:"user"`
	Message string       `json:"message,omitempty"`
}

func (h *Handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := req.Context()
	var reqBody SignupRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	_, err := h.users.GetByEmail(ctx, reqBody.Email)
	if err == nil {
		http.Error(w, "user with this email already exists", http.StatusConflict)
		return
	}

	// Check if username already exists
	_, err = h.users.GetByUsername(ctx, reqBody.Username)
	if err == nil {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	// Create user
	user := &models.User{
		Username: reqBody.Username,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	if err := h.users.Create(ctx, user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, AuthResponse{
		Token:   token,
		User:    user,
		Message: "signup successful",
	}, http.StatusCreated)
}

func (h *Handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := req.Context()
	var loginReq LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by email
	user, err := h.users.GetByEmail(ctx, loginReq.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := h.users.VerifyPassword(ctx, user.ID, loginReq.Password); err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, AuthResponse{
		Token:   token,
		User:    user,
		Message: "login successful",
	}, http.StatusOK)
}

func (h *Handler) logoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a JWT-based system, logout is typically handled on the client-side by deleting the token.
	// This endpoint is provided as a convenience for API documentation.
	writeJSON(w, map[string]string{
		"message": "logout successful",
	}, http.StatusOK)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/brennanromance/heard/internal/models"
	"github.com/brennanromance/heard/internal/repo"
)

type Handler struct {
	companies *repo.CompanyRepo
	users     *repo.UserRepo
	posts     *repo.PostRepo
	comments  *repo.CommentRepo
}

func NewHandler(c *repo.CompanyRepo, u *repo.UserRepo, p *repo.PostRepo, cm *repo.CommentRepo) *Handler {
	return &Handler{companies: c, users: u, posts: p, comments: cm}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Auth routes (no middleware)
	mux.HandleFunc("POST /signup", h.signupHandler)
	mux.HandleFunc("POST /login", h.loginHandler)
	mux.HandleFunc("POST /logout", h.logoutHandler)

	// Protected routes
	mux.HandleFunc("GET /companies", h.AuthMiddleware(h.companiesHandlerGET))
	mux.HandleFunc("POST /companies", h.AuthMiddleware(h.companiesHandlerPOST))
	mux.HandleFunc("PUT /companies", h.AuthMiddleware(h.companiesHandlerPUT))
	mux.HandleFunc("DELETE /companies", h.AuthMiddleware(h.companiesHandlerDELETE))

	mux.HandleFunc("GET /posts", h.AuthMiddleware(h.postsHandlerGET))
	mux.HandleFunc("POST /posts", h.AuthMiddleware(h.postsHandlerPOST))
	mux.HandleFunc("PUT /posts", h.AuthMiddleware(h.postsHandlerPUT))
	mux.HandleFunc("DELETE /posts", h.AuthMiddleware(h.postsHandlerDELETE))

	mux.HandleFunc("GET /comments", h.AuthMiddleware(h.commentsHandlerGET))
	mux.HandleFunc("POST /comments", h.AuthMiddleware(h.commentsHandlerPOST))
	mux.HandleFunc("PUT /comments", h.AuthMiddleware(h.commentsHandlerPUT))
	mux.HandleFunc("DELETE /comments", h.AuthMiddleware(h.commentsHandlerDELETE))
}

func idFromQuery(r *http.Request) (int, bool) {
	qs := r.URL.Query().Get("id")
	if qs == "" {
		return 0, false
	}
	id, err := strconv.Atoi(qs)
	if err != nil {
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// Auth Handlers

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

func (h *Handler) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	_, err := h.users.GetByEmail(ctx, req.Email)
	if err == nil {
		http.Error(w, "user with this email already exists", http.StatusConflict)
		return
	}

	// Check if username already exists
	_, err = h.users.GetByUsername(ctx, req.Username)
	if err == nil {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
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

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by email
	user, err := h.users.GetByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := h.users.VerifyPassword(ctx, user.ID, req.Password); err != nil {
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

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a JWT-based system, logout is typically handled on the client-side by deleting the token.
	// This endpoint is provided as a convenience for API documentation.
	writeJSON(w, map[string]string{
		"message": "logout successful",
	}, http.StatusOK)
}

// Companies

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

// Posts

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

// Comments

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

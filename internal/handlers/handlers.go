package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

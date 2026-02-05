package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/brennanromance/heard/internal/models"
    "github.com/brennanromance/heard/internal/repo"
)

type Handler struct{
    companies *repo.CompanyRepo
    users     *repo.UserRepo
    posts     *repo.PostRepo
    comments  *repo.CommentRepo
}

func NewHandler(c *repo.CompanyRepo, u *repo.UserRepo, p *repo.PostRepo, cm *repo.CommentRepo) *Handler {
    return &Handler{companies: c, users: u, posts: p, comments: cm}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/companies", h.companiesHandler)
    mux.HandleFunc("/users", h.usersHandler)
    mux.HandleFunc("/posts", h.postsHandler)
    mux.HandleFunc("/comments", h.commentsHandler)
}

func idFromQuery(r *http.Request) (int, bool) {
    qs := r.URL.Query().Get("id")
    if qs == "" { return 0, false }
    id, err := strconv.Atoi(qs)
    if err != nil { return 0, false }
    return id, true
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}

// Companies
func (h *Handler) companiesHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    switch r.Method {
    case http.MethodGet:
        if id, ok := idFromQuery(r); ok {
            c, err := h.companies.GetByID(ctx, id)
            if err != nil { http.Error(w, err.Error(), http.StatusNotFound); return }
            writeJSON(w, c, http.StatusOK); return
        }
        list, err := h.companies.List(ctx)
        if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, list, http.StatusOK)
    case http.MethodPost:
        var c models.Company
        if err := json.NewDecoder(r.Body).Decode(&c); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        if err := h.companies.Create(ctx, &c); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, c, http.StatusCreated)
    case http.MethodPut:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        var c models.Company
        if err := json.NewDecoder(r.Body).Decode(&c); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        c.ID = id
        if err := h.companies.Update(ctx, &c); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, c, http.StatusOK)
    case http.MethodDelete:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        if err := h.companies.Delete(ctx, id); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

// Users
func (h *Handler) usersHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    switch r.Method {
    case http.MethodGet:
        if id, ok := idFromQuery(r); ok {
            u, err := h.users.GetByID(ctx, id)
            if err != nil { http.Error(w, err.Error(), http.StatusNotFound); return }
            writeJSON(w, u, http.StatusOK); return
        }
        list, err := h.users.List(ctx)
        if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, list, http.StatusOK)
    case http.MethodPost:
        var u models.User
        if err := json.NewDecoder(r.Body).Decode(&u); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        if err := h.users.Create(ctx, &u); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, u, http.StatusCreated)
    case http.MethodPut:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        var u models.User
        if err := json.NewDecoder(r.Body).Decode(&u); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        u.ID = id
        if err := h.users.Update(ctx, &u); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, u, http.StatusOK)
    case http.MethodDelete:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        if err := h.users.Delete(ctx, id); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

// Posts
func (h *Handler) postsHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    switch r.Method {
    case http.MethodGet:
        if id, ok := idFromQuery(r); ok {
            p, err := h.posts.GetByID(ctx, id)
            if err != nil { http.Error(w, err.Error(), http.StatusNotFound); return }
            writeJSON(w, p, http.StatusOK); return
        }
        list, err := h.posts.List(ctx)
        if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, list, http.StatusOK)
    case http.MethodPost:
        var p models.Post
        if err := json.NewDecoder(r.Body).Decode(&p); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        if err := h.posts.Create(ctx, &p); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, p, http.StatusCreated)
    case http.MethodPut:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        var p models.Post
        if err := json.NewDecoder(r.Body).Decode(&p); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        p.ID = id
        if err := h.posts.Update(ctx, &p); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, p, http.StatusOK)
    case http.MethodDelete:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        if err := h.posts.Delete(ctx, id); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

// Comments
func (h *Handler) commentsHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    switch r.Method {
    case http.MethodGet:
        if id, ok := idFromQuery(r); ok {
            c, err := h.comments.GetByID(ctx, id)
            if err != nil { http.Error(w, err.Error(), http.StatusNotFound); return }
            writeJSON(w, c, http.StatusOK); return
        }
        list, err := h.comments.List(ctx)
        if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, list, http.StatusOK)
    case http.MethodPost:
        var c models.Comment
        if err := json.NewDecoder(r.Body).Decode(&c); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        if err := h.comments.Create(ctx, &c); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, c, http.StatusCreated)
    case http.MethodPut:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        var c models.Comment
        if err := json.NewDecoder(r.Body).Decode(&c); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        c.ID = id
        if err := h.comments.Update(ctx, &c); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        writeJSON(w, c, http.StatusOK)
    case http.MethodDelete:
        id, ok := idFromQuery(r)
        if !ok { http.Error(w, "missing id", http.StatusBadRequest); return }
        if err := h.comments.Delete(ctx, id); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

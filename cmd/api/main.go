package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/brennanromance/heard/internal/db"
    "github.com/brennanromance/heard/internal/handlers"
    "github.com/brennanromance/heard/internal/repo"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://localhost:5432/heard?sslmode=disable"
    }

    sqlDB, err := db.Connect(context.Background(), dsn)
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }
    defer sqlDB.Close()

    // repositories
    companyRepo := repo.NewCompanyRepo(sqlDB)
    userRepo := repo.NewUserRepo(sqlDB)
    postRepo := repo.NewPostRepo(sqlDB)
    commentRepo := repo.NewCommentRepo(sqlDB)

    // handlers
    h := handlers.NewHandler(companyRepo, userRepo, postRepo, commentRepo)

    mux := http.NewServeMux()
    h.RegisterRoutes(mux)

    addr := ":8080"
    log.Printf("listening on %s", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatalf("server: %v", err)
    }
}

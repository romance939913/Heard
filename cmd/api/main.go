package main

import (
    "context"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/brennanromance/heard/internal/db"
    "github.com/brennanromance/heard/internal/handlers"
    "github.com/brennanromance/heard/internal/repo"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
        log.Fatalf("DB_USER, DB_PASSWORD, DB_HOST, or DB_PORT environment variables are not set")
    }

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://" + url.QueryEscape(dbUser) + ":" + url.QueryEscape(dbPassword) +
            "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
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
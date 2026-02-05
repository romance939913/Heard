package db

import (
    "context"
    "database/sql"

    _ "github.com/jackc/pgx/v5/stdlib"
)

// Connect opens a *sql.DB using the provided connection string via pgx stdlib.
func Connect(ctx context.Context, connString string) (*sql.DB, error) {
    db, err := sql.Open("pgx", connString)
    if err != nil {
        return nil, err
    }
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}

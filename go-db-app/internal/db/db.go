package db

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver
)

// Connect establishes a connection to the database and returns a *sql.DB instance.
func Connect(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
        return nil, err
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
        return nil, err
    }

    return db, nil
}
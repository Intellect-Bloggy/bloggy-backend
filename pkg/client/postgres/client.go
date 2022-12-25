package postgres

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
    Host     string
    Port     uint16
    Username string
    Password string
    DBName   string
    SSLMode  string
}

func NewPostgresDB(cfg PostgresConfig) (*pgx.Conn, error) {

    connStr := fmt.Sprintf(`host=%s port=%d user=%s dbname=%s password=%s sslmode=%s`,
        cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
    db, err := pgx.Connect(context.Background(), connStr)

    if err != nil {
        return nil, err
    }

    err = db.Ping(context.Background())
    if err != nil {
        return nil, err
    }

    return db, nil
}

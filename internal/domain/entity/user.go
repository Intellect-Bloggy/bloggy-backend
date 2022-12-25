package entity

import (
    "time"

    "github.com/jackc/pgx/v5/pgtype"
)

type User struct {
    Id        int         `db:"id"`
    Username  string      `db:"username"`
    Password  string      `db:"password"`
    Name      pgtype.Text `db:"name"`
    Birthday  pgtype.Date `db:"birthday"`
    Email     pgtype.Text `db:"email"`
    Phone     pgtype.Text `db:"phone"`
    CreatedAt time.Time   `db:"created_at"`
}

type UserResponse struct {
    Username  string      `json:"username"`
    Name      pgtype.Text `json:"name"`
    Birthday  pgtype.Date `json:"birthday"`
    Email     pgtype.Text `json:"email"`
    Phone     pgtype.Text `json:"phone"`
    CreatedAt time.Time   `json:"created_at"`
}

type UserResponseSwagger struct {
    Username  string    `json:"username"`
    Name      string    `json:"name"`
    Birthday  string    `json:"birthday"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
}

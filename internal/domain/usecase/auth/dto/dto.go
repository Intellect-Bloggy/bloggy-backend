package dto

import (
    "github.com/jackc/pgx/v5/pgtype"
)

type SignUpDTO struct {
    Username string      `db:"username"`
    Password string      `db:"password"`
    Name     pgtype.Text `db:"name"`
    Birthday pgtype.Date `db:"birthday"`
    Email    pgtype.Text `db:"email"`
    Phone    pgtype.Text `db:"phone"`
}

type LogoutDTO struct {
    RefreshToken string `json:"token" binding:"required"`
}

type RefreshDTO struct {
    RefreshToken string `json:"token" binding:"required"`
}

package entity

import (
    "time"
)

type Auth struct {
    Access  string `json:"access_token"`
    Refresh string `json:"refresh_token"`
}

type Session struct {
    RefreshToken string    `db:"token"`
    ExpiresAt    time.Time `db:"expires_at"`
}

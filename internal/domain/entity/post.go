package entity

import (
    "time"
)

type Post struct {
    Id        int       `json:"id" db:"id"`
    AuthorId  int       `json:"author_id" db:"author_id"`
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreatePostResponse struct {
    Id int `json:"id"`
}

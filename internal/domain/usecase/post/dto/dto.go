package dto

type CreatePostDTO struct {
    Content string `json:"content" binding:"required"`
}

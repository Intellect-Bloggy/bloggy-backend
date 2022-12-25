package dto

type CreateUserDTO struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`

    Name     *string `json:"name"`
    Birthday *string `json:"birthday"` // Use "2000-12-31" format
    Email    *string `json:"email"`
    Phone    *string `json:"phone"`
}

type GetByCredentialsDTO struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type EditUserDTO struct {
    Username *string `json:"username"`
    Name     *string `json:"name"`
    Birthday *string `json:"birthday"`
    Email    *string `json:"email"`
    Phone    *string `json:"phone"`
}

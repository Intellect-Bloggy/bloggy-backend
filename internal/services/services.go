package services

import (
    "context"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/hash"
)

type Services struct {
    User
    Auth
}

type User interface {
    GetUserByUsername(ctx context.Context, username string) (UserResponse, error)
    EditById(ctx context.Context, id int, input EditInput) (UserResponse, error)
}

type Auth interface {
    SignUp(ctx context.Context, input SignUpInput) (Tokens, error)
    SignIn(ctx context.Context, input SignInInput) (Tokens, error)
    RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

func NewServices(repos *repository.Repository, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *Services {
    return &Services{
        User: newUserService(repos),
        Auth: newAuthService(repos, hasher, tokenManager),
    }
}

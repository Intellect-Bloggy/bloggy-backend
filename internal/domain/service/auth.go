package service

import (
    "context"
    "strconv"
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    auth_helpers "github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
)

type authStorage interface {
    SetSession(ctx context.Context, userId int, session entity.Session) error
    DeleteRefresh(ctx context.Context, refreshToken string) error
    CheckRefresh(ctx context.Context, refreshToken string) error
}

type authService struct {
    storage      authStorage
    tokenManager auth_helpers.TokenManager
}

func NewAuthService(storage authStorage, tokenManager auth_helpers.TokenManager) *authService {
    return &authService{storage: storage, tokenManager: tokenManager}
}

func (s *authService) GenerateAccessToken(ctx context.Context, userId int) (string, error) {
    return s.tokenManager.NewJWT(strconv.Itoa(userId), time.Minute*15)
}

func (s *authService) GenerateRefreshToken(ctx context.Context) (string, error) {
    return s.tokenManager.NewRefreshToken()
}

func (s *authService) SetSession(ctx context.Context, userId int, session entity.Session) error {
    return s.storage.SetSession(ctx, userId, session)
}

func (s *authService) CheckRefresh(ctx context.Context, refreshToken string) error {
    return s.storage.CheckRefresh(ctx, refreshToken)
}

func (s *authService) DeleteRefresh(ctx context.Context, refreshToken string) error {
    return s.storage.DeleteRefresh(ctx, refreshToken)
}

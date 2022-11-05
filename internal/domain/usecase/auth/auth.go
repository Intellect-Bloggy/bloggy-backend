package auth

import (
    "context"
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/auth/dto"
    user_usecase "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/user/dto"
)

type Service interface {
    GenerateAccessToken(ctx context.Context, userId int) (string, error)
    GenerateRefreshToken(ctx context.Context) (string, error)
    SetSession(ctx context.Context, userId int, session entity.Session) error
    CheckRefresh(ctx context.Context, refreshToken string) error
    DeleteRefresh(ctx context.Context, refreshToken string) error
}

type UserService interface {
    CreateUser(ctx context.Context, dto user_usecase.CreateUserDTO) (int, error)
    GetUserByID(ctx context.Context, id int) (entity.User, error)
    GetByCredentials(ctx context.Context, dto user_usecase.GetByCredentialsDTO) (entity.User, error)
    GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
}

type authUsecase struct {
    Service
    UserService
}

func NewAuthUsecase(s Service, us UserService) *authUsecase {
    return &authUsecase{
        Service:     s,
        UserService: us,
    }
}

func (u *authUsecase) SignUp(ctx context.Context, dto user_usecase.CreateUserDTO) (entity.Auth, error) {

    response := entity.Auth{}

    newUserID, err := u.UserService.CreateUser(ctx, dto)
    if err != nil {
        return response, err
    }

    user, err := u.GetUserByID(ctx, newUserID)
    if err != nil {
        return response, err
    }

    response.Access, err = u.GenerateAccessToken(ctx, user.Id)
    if err != nil {
        return response, err
    }

    response.Refresh, err = u.GenerateRefreshToken(ctx)
    if err != nil {
        return response, err
    }

    session := entity.Session{
        RefreshToken: response.Refresh,
        ExpiresAt:    time.Now().Add(time.Hour * 720), // 30 days
    }

    err = u.SetSession(ctx, user.Id, session)

    return response, err
}

func (u *authUsecase) SignIn(ctx context.Context, dto user_usecase.GetByCredentialsDTO) (entity.Auth, error) {
    response := entity.Auth{}
    user, err := u.UserService.GetByCredentials(ctx, dto)
    if err != nil {
        return response, err
    }

    response.Access, err = u.GenerateAccessToken(ctx, user.Id)
    if err != nil {
        return response, err
    }

    response.Refresh, err = u.GenerateRefreshToken(ctx)
    if err != nil {
        return response, err
    }

    session := entity.Session{
        RefreshToken: response.Refresh,
        ExpiresAt:    time.Now().Add(time.Hour * 720), // 30 days
    }

    err = u.SetSession(ctx, user.Id, session)

    return response, err
}

func (u *authUsecase) Refresh(ctx context.Context, dto dto.RefreshDTO) (entity.Auth, error) {
    response := entity.Auth{}
    var err error
    if err = u.Service.CheckRefresh(ctx, dto.RefreshToken); err != nil {
        if err = u.Service.DeleteRefresh(ctx, dto.RefreshToken); err != nil {
            return response, err
        }
    }

    var user entity.User
    user, err = u.UserService.GetByRefreshToken(ctx, dto.RefreshToken)

    response.Access, err = u.GenerateAccessToken(ctx, user.Id)
    if err != nil {
        return response, err
    }

    response.Refresh, err = u.GenerateRefreshToken(ctx)
    if err != nil {
        return response, err
    }

    session := entity.Session{
        RefreshToken: response.Refresh,
        ExpiresAt:    time.Now().Add(time.Hour * 720), // 30 days
    }

    err = u.SetSession(ctx, user.Id, session)

    return response, err
}

func (u *authUsecase) Logout(ctx context.Context, dto dto.LogoutDTO) error {
    if err := u.Service.CheckRefresh(ctx, dto.RefreshToken); err != nil {
        return err
    }

    return u.Service.DeleteRefresh(ctx, dto.RefreshToken)
}
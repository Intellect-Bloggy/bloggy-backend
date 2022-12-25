package service

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgtype"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    user_dto "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/user/dto"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/hash"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/utils"
)

type userStorage interface {
    GetUserByID(ctx context.Context, id int) (entity.User, error)
    GetByUsername(ctx context.Context, username string) (entity.UserResponse, error)
    CreateUser(ctx context.Context, user entity.User) (int, error)
    GetUserByUserNameAndPassword(ctx context.Context, username, password string) (entity.User, error)
    GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
    EditById(ctx context.Context, id int, req map[string]string) (entity.UserResponse, error)
}

type userService struct {
    storage userStorage
    hasher  hash.PasswordHasher
}

func NewUserService(storage userStorage, hasher hash.PasswordHasher) *userService {
    return &userService{storage: storage, hasher: hasher}
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
    return s.storage.GetUserByID(ctx, id)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (entity.UserResponse, error) {
    return s.storage.GetByUsername(ctx, username)
}

func (s *userService) CreateUser(ctx context.Context, dto user_dto.CreateUserDTO) (int, error) {

    if err := utils.CheckUsername(dto.Username); err != nil {
        return -1, err
    }

    if err := utils.CheckPassword(dto.Password); err != nil {
        return -1, err
    }

    var name pgtype.Text
    if dto.Name != nil && *dto.Name != "" {
        name = pgtype.Text{String: *dto.Name, Valid: true}
    }

    var birthday pgtype.Date
    if dto.Birthday != nil && *dto.Birthday != "" {
        date, err := utils.ParseDate(*dto.Birthday)
        if err != nil {
            return -1, err
        }
        birthday = pgtype.Date{Time: date, Valid: true}
    }

    var email pgtype.Text
    if dto.Email != nil && *dto.Email != "" {
        email = pgtype.Text{String: *dto.Email, Valid: true}
    }

    var phone pgtype.Text
    if dto.Phone != nil && *dto.Phone != "" {
        phone = pgtype.Text{String: *dto.Phone, Valid: true}
    }

    passHash, err := s.hasher.Hash(dto.Password)
    if err != nil {
        return -1, err

    }

    u := entity.User{
        Username:  dto.Username,
        Password:  passHash,
        Name:      name,
        Birthday:  birthday,
        Email:     email,
        Phone:     phone,
        CreatedAt: time.Now(),
    }

    return s.storage.CreateUser(ctx, u)
}

func (s *userService) GetByCredentials(ctx context.Context, dto user_dto.GetByCredentialsDTO) (entity.User, error) {
    passHash, err := s.hasher.Hash(dto.Password)
    if err != nil {
        return entity.User{}, err
    }

    return s.storage.GetUserByUserNameAndPassword(ctx, dto.Username, passHash)
}

func (s *userService) GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {
    return s.storage.GetByRefreshToken(ctx, refreshToken)
}

func (s *userService) EditById(ctx context.Context, id int, i user_dto.EditUserDTO) (entity.UserResponse, error) {
    if id <= 0 {
        return entity.UserResponse{}, errors.ErrWrongId
    }

    reqMap := map[string]string{}
    if i.Username != nil {
        if err := utils.CheckUsername(*i.Username); err != nil {
            return entity.UserResponse{}, err
        }
        reqMap["username"] = *i.Username
    }
    if i.Name != nil && *i.Name != "" {
        if err := utils.CheckName(*i.Name); err != nil {
            return entity.UserResponse{}, err
        }
        reqMap["name"] = *i.Name
    }
    if i.Birthday != nil && *i.Birthday != "" {
        // TODO: Добавить проверку дня рождения на совпадение с форматом 2000-12-31
        reqMap["birthday"] = *i.Birthday
    }
    if i.Email != nil && *i.Email != "" {
        // TODO: Добавить проверку почты на совпадение с форматом xxx@zzz.com
        reqMap["email"] = *i.Email
    }
    if i.Phone != nil && *i.Phone != "" {
        // TODO: Сделать валидацию номера
        reqMap["phone"] = *i.Phone
    }

    user, err := s.storage.EditById(ctx, id, reqMap)
    if err != nil {
        return entity.UserResponse{}, err
    }

    return user, nil
}

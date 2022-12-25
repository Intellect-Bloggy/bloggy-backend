package user

import (
    "context"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    user_dto "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/user/dto"
)

type service interface {
    GetByUsername(ctx context.Context, username string) (entity.UserResponse, error)
    EditById(ctx context.Context, id int, i user_dto.EditUserDTO) (entity.UserResponse, error)
}

type postService interface {
    GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error)
}

type userUsecase struct {
    service
    postService
}

func NewUserUsecase(service service, postService postService) *userUsecase {
    return &userUsecase{service: service, postService: postService}
}

func (u *userUsecase) GetByUsername(ctx context.Context, username string) (entity.UserResponse, error) {
    return u.service.GetByUsername(ctx, username)
}

func (u *userUsecase) GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error) {
    return u.postService.GetAllByUsername(ctx, username)
}

func (u *userUsecase) EditById(ctx context.Context, id int, dto user_dto.EditUserDTO) (entity.UserResponse, error) {
    return u.service.EditById(ctx, id, dto)
}

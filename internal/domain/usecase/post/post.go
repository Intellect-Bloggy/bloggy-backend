package post

import (
    "context"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/post/dto"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type postService interface {
    Create(ctx context.Context, p entity.Post) (id int, err error)
    GetById(ctx context.Context, id int) (p entity.Post, err error)
    DeleteById(ctx context.Context, id int) error
    IsAuthor(c context.Context, postId int, authorId int) (bool, error)
}

type postUsecase struct {
    postService
}

func NewPostUsecase(postService postService) *postUsecase {
    return &postUsecase{postService: postService}
}

func (u *postUsecase) Create(ctx context.Context, dto dto.CreatePostDTO, userId int) (int, error) {
    post := entity.Post{}
    post.AuthorId = userId
    post.Content = dto.Content

    if len(dto.Content) < 1 {
        return -1, errors.ErrPostContent
    }

    return u.postService.Create(ctx, post)
}

func (u *postUsecase) GetById(ctx context.Context, id int) (entity.Post, error) {
    return u.postService.GetById(ctx, id)
}

func (u *postUsecase) DeleteById(ctx context.Context, id int, userId int) error {
    if ok, err := u.IsAuthor(ctx, id, userId); err != nil || !ok {
        if err != nil {
            return err
        }
        if !ok {
            return errors.ErrUserIsNotAuthor
        }
    }

    if err := u.postService.DeleteById(ctx, id); err != nil {
        return err
    }

    return nil
}

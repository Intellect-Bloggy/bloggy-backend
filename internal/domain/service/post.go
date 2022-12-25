package service

import (
    "context"
    "time"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
)

type postStorage interface {
    Create(ctx context.Context, p entity.Post) (id int, err error)
    GetById(ctx context.Context, id int) (p entity.Post, err error)
    DeleteById(ctx context.Context, id int) error
    GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error)
    IsAuthor(c context.Context, postId int, authorId int) (bool, error)
}

type postService struct {
    postStorage
}

func NewPostService(postStorage postStorage) *postService {
    return &postService{postStorage: postStorage}
}

func (s *postService) Create(ctx context.Context, p entity.Post) (id int, err error) {
    p.CreatedAt = time.Now()
    return s.postStorage.Create(ctx, p)
}

func (s *postService) GetById(ctx context.Context, id int) (p entity.Post, err error) {
    return s.postStorage.GetById(ctx, id)
}

func (s *postService) DeleteById(ctx context.Context, id int) error {
    return s.postStorage.DeleteById(ctx, id)
}

func (s *postService) GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error) {
    return s.postStorage.GetAllByUsername(ctx, username)
}

func (s *postService) IsAuthor(ctx context.Context, postId int, authorId int) (bool, error) {
    return s.postStorage.IsAuthor(ctx, postId, authorId)
}

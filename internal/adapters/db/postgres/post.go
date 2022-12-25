package postgres

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type postStorage struct {
    db *pgx.Conn
}

func NewPostStorage(db *pgx.Conn) *postStorage {
    return &postStorage{db: db}
}

func (s *postStorage) Create(ctx context.Context, p entity.Post) (id int, err error) {
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        INSERT INTO %s(author_id, content, created_at)
        VALUES ($1, $2, $3)
        RETURNING id
    `, postsTable), p.AuthorId, p.Content, p.CreatedAt).Scan(&id)
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (s *postStorage) GetById(ctx context.Context, id int) (p entity.Post, err error) {
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id, author_id, content, created_at
        FROM %s
        WHERE id = $1
    `, postsTable), id).Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return entity.Post{}, e.ErrIdNotFound
        }

        return entity.Post{}, err
    }

    return p, nil
}

func (s *postStorage) DeleteById(ctx context.Context, id int) error {
    var postId int
    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        DELETE FROM %s
        WHERE id = $1
        RETURNING id
    `, postsTable), id).Scan(&postId)
    if errors.Is(err, pgx.ErrNoRows) {
        return e.ErrIdNotFound
    }

    return err
}

func (s *postStorage) GetAllByUsername(ctx context.Context, username string) (posts []entity.Post, err error) {
    var authorId int
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), username).Scan(&authorId)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return make([]entity.Post, 0), e.ErrUsernameNotFound
        }

        return make([]entity.Post, 0), err
    }

    rows, err := s.db.Query(ctx, fmt.Sprintf(`
        SELECT id, author_id, content, created_at 
        FROM %s WHERE author_id = $1
    `, postsTable), authorId)
    if err != nil {
        return make([]entity.Post, 0), err
    }

    posts = make([]entity.Post, 0)
    for rows.Next() {
        var p entity.Post
        err = rows.Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
        if err != nil {
            return make([]entity.Post, 0), err
        }
        posts = append(posts, p)
    }

    return posts, nil
}

func (s *postStorage) IsAuthor(c context.Context, postId int, authorId int) (bool, error) {
    p := entity.Post{}
    err := s.db.QueryRow(c, fmt.Sprintf(`
        SELECT id, author_id, content, created_at
        FROM %s
        WHERE id = $1 AND author_id = $2
    `, postsTable), postId, authorId).Scan(&p.Id, &p.AuthorId, &p.Content, &p.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return false, nil
        }

        return false, err
    }

    return true, nil
}

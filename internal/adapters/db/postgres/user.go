package postgres

import (
    "context"
    "errors"
    "fmt"
    "strings"

    "github.com/jackc/pgx/v5"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    e "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type userStorage struct {
    db *pgx.Conn
}

func NewUserStorage(db *pgx.Conn) *userStorage {
    return &userStorage{db: db}
}

func (s *userStorage) GetUserByID(ctx context.Context, id int) (entity.User, error) {
    user := entity.User{}
    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id, username, name, birthday, email, phone, created_at
        FROM %s
        WHERE id = $1
    `, usersTable), id).Scan(
        &user.Id, &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return entity.User{}, e.ErrIdNotFound
    }
    if err != nil {
        return entity.User{}, err
    }

    return user, nil
}

func (s *userStorage) GetByUsername(ctx context.Context, username string) (entity.UserResponse, error) {
    user := entity.UserResponse{}
    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT username, name, birthday, email, phone, created_at
        FROM %s
        WHERE username = $1
    `, usersTable), username).Scan(
        &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return entity.UserResponse{}, e.ErrIdNotFound
    }
    if err != nil {
        return entity.UserResponse{}, err
    }

    return user, nil
}

func (s *userStorage) CreateUser(ctx context.Context, u entity.User) (int, error) {
    _id := 0
    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), u.Username).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return -1, e.ErrTakenUsername
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return -1, err
    }

    tx, err := s.db.Begin(ctx)

    newU := entity.User{}

    // FIXME: Тут нет проверки занятости email и phone
    // Создание пользователя
    err = tx.QueryRow(ctx, fmt.Sprintf(`
        INSERT INTO %s (username, name, email, phone, birthday, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, username, name, email, phone, birthday, created_at
    `, usersTable), u.Username, u.Name, u.Email, u.Phone, u.Birthday, u.CreatedAt).Scan(
        &newU.Id, &newU.Username, &newU.Name, &newU.Email, &newU.Phone, &newU.Birthday, &newU.CreatedAt)
    if err != nil {
        _ = tx.Rollback(ctx)
        return -1, err
    }

    // Привязка пароля
    err = tx.QueryRow(ctx, fmt.Sprintf(`
        INSERT INTO %s (user_id, password)
        VALUES ($1, $2)
        RETURNING user_id, password
    `, authTable), newU.Id, u.Password).Scan(&newU.Id, &newU.Password)
    if err != nil {
        _ = tx.Rollback(ctx)
        return -1, err
    }

    err = tx.Commit(ctx)
    if err != nil {
        return -1, err
    }

    return newU.Id, nil
}

func (s *userStorage) GetUserByUserNameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
    user := entity.User{}

    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT u.id, u.username, u.name, u.birthday, u.email, u.phone, u.created_at
        FROM %s AS a
        JOIN %s AS u
        ON u.id = a.user_id
        WHERE a.password = $1 AND u.username = $2
    `, authTable, usersTable), password, username).Scan(
        &user.Id, &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return entity.User{}, e.ErrWrongCredentials
    }
    if err != nil {
        return entity.User{}, err
    }

    return user, nil
}

func (s *userStorage) GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {
    user := entity.User{}

    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT u.id, u.username, u.name, u.birthday, u.email, u.phone, u.created_at
        FROM %s AS r
        JOIN %s AS u
        ON u.id = r.user_id
        WHERE r.token = $1
    `, refreshTable, usersTable), refreshToken).Scan(
        &user.Id, &user.Username, &user.Name, &user.Birthday,
        &user.Email, &user.Phone, &user.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return entity.User{}, e.ErrUserNotFound
    }
    if err != nil {
        return entity.User{}, err
    }

    return user, nil
}

func (s *userStorage) EditById(ctx context.Context, id int, req map[string]string) (entity.UserResponse, error) {
    // Проверяем доступность username
    _id := 0
    err := s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE username = $1
    `, usersTable), req["username"]).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return entity.UserResponse{}, e.ErrTakenUsername
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return entity.UserResponse{}, err
    }

    // Проверяем доступность email
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE email = $1
    `, usersTable), req["email"]).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return entity.UserResponse{}, e.ErrTakenEmail
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return entity.UserResponse{}, err
    }

    // Проверяем доступность phone
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        SELECT id
        FROM %s
        WHERE phone = $1
    `, usersTable), req["phone"]).Scan(&_id)
    if err == nil {
        // Если он нашел пользователя и успешно просканировал, то он существует
        return entity.UserResponse{}, e.ErrTakenPhone
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        // Если ошибка не является ошибкой "Не найден пользователь"
        return entity.UserResponse{}, err
    }

    // Создаем параметры
    var updateQuery []string
    for key, value := range req {
        updateQuery = append(updateQuery, fmt.Sprintf(`%s = '%s'`, key, value))
    }

    user := entity.UserResponse{}
    err = s.db.QueryRow(ctx, fmt.Sprintf(`
        UPDATE %s 
        SET %s
        WHERE id = $1
        RETURNING username, name, phone, email, birthday
    `, usersTable, strings.Join(updateQuery, ", ")), id).Scan(
        &user.Username, &user.Name, &user.Phone, &user.Email, &user.Birthday)
    if err != nil {
        return entity.UserResponse{}, err
    }

    return user, nil
}

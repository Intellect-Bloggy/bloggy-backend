CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(60),
    username   VARCHAR(30) NOT NULL UNIQUE,

    email      VARCHAR(50) UNIQUE,
    phone      VARCHAR(15) UNIQUE,
    birthday   DATE,
    created_at TIMESTAMP   NOT NULL
);

CREATE TABLE posts
(
    id         SERIAL PRIMARY KEY,
    author_id  INTEGER   NOT NULL,

    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE auth
(
    user_id  INTEGER UNIQUE NOT NULL,
    password VARCHAR(255)   NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE refresh
(
    user_id    INTEGER      NOT NULL,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP    NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

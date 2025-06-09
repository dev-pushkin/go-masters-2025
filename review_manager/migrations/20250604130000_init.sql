-- +goose Up

-- Таблица пользователей
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    email      TEXT UNIQUE NOT NULL
);

-- Таблица объектов
CREATE TABLE objects (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    name       TEXT NOT NULL
);


-- Таблица отзывов
CREATE TABLE reviews (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    status      TEXT NOT NULL CHECK (status != ''),
    text        TEXT NOT NULL,
    rate_level  INTEGER NOT NULL CHECK (rate_level >= 0 AND rate_level <= 5),
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_id   BIGINT NOT NULL REFERENCES objects(id) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_object_id ON reviews(object_id);
CREATE INDEX idx_reviews_status ON reviews(status);
CREATE INDEX idx_reviews_created_at ON reviews(created_at);

-- +goose Down

DROP INDEX IF EXISTS idx_reviews_created_at;
DROP INDEX IF EXISTS idx_reviews_status;
DROP INDEX IF EXISTS idx_reviews_object_id;
DROP INDEX IF EXISTS idx_reviews_user_id;

DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS objects;
DROP TABLE IF EXISTS users;

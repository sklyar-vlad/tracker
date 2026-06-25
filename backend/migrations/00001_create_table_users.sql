-- +goose Up
CREATE TYPE ROLE AS ENUM('user', 'admin') ;

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role ROLE NOT NULL DEFAULT 'user', 
    username VARCHAR(30) UNIQUE,
    email VARCHAR(100) UNIQUE NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS ROLE;
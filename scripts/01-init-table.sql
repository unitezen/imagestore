\c imagestore

CREATE TABLE users (
    id serial PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    password VARCHAR(64) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE user_sessions (
    api_key VARCHAR PRIMARY KEY,
    user_id INTEGER REFERENCES users(id)
);

CREATE TABLE images (
    id serial PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name    VARCHAR(255) NOT NULL,
    data  VARCHAR NOT NULL
);
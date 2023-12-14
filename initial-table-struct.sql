-- SQL to create the users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255)
);

-- SQL to create the tokens table if it doesn't exist
CREATE TABLE IF NOT EXISTS user_token (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    token VARCHAR(255),
    expired BOOLEAN
);

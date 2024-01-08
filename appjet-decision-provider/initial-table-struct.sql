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

-- SQL to create the appjet_configurations table if it doesn't exist
CREATE TABLE IF NOT EXISTS appjet_configurations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    config TEXT
);

-- SQL to create the root user
INSERT INTO users (username, email, password)
VALUES ('root', 'admin@admin.com', 'admin');

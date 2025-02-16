CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    coins INTEGER NOT NULL DEFAULT 1000
    );

CREATE TABLE IF NOT EXISTS purchases (
                                         id SERIAL PRIMARY KEY,
                                         user_id INTEGER REFERENCES users(id),
    item VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    UNIQUE (user_id, item)
    );

CREATE TABLE IF NOT EXISTS transactions (
                                            id SERIAL PRIMARY KEY,
                                            user_id INTEGER REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    counterpart VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE TABLE IF NOT EXISTS users (
    id      BIGINT PRIMARY KEY,
    balance NUMERIC NOT NULL CHECK (balance >= 0)
);

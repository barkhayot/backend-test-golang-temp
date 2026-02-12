CREATE TABLE IF NOT EXISTS user_balance_history (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL,
    balance_from NUMERIC NOT NULL,
    balance_to   NUMERIC NOT NULL,
    amount       NUMERIC NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

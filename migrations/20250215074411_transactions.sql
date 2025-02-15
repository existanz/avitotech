-- +goose Up
-- +goose StatementBegin
CREATE TABLE coin_transactions (
                                   id SERIAL PRIMARY KEY,
                                   from_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                                   to_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                                   amount INTEGER NOT NULL,
                                   transaction_type VARCHAR(50) NOT NULL,
                                   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE coin_transactions;
-- +goose StatementEnd

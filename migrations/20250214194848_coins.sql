-- +goose Up
-- +goose StatementBegin
CREATE TABLE coins (
                       user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       amount INTEGER NOT NULL DEFAULT 0,
                       PRIMARY KEY (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE coins;
-- +goose StatementEnd

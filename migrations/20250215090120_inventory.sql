-- +goose Up
-- +goose StatementBegin
CREATE TABLE inventory (
                           user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           item_type VARCHAR(255) NOT NULL,
                           quantity INTEGER NOT NULL DEFAULT 0,
                           PRIMARY KEY (user_id, item_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inventory;
-- +goose StatementEnd

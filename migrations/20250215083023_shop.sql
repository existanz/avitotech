-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop (
                    id SERIAL PRIMARY KEY,
                    item_type varchar(255) NOT NULL,
                    price INTEGER NOT NULL
);

CREATE INDEX idx_shop_item_type ON shop(item_type);

INSERT INTO shop (item_type, price) VALUES
                                        ( 't-shirt', 80),
                                        ( 'cup', 20),
                                        ( 'book', 50),
                                        ( 'pen', 10),
                                        ( 'powerbank', 200),
                                        ( 'hoody', 300),
                                        ( 'umbrella', 200),
                                        ( 'socks', 10),
                                        ( 'wallet', 50),
                                        ( 'pink-hoody', 500);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop;
-- +goose StatementEnd

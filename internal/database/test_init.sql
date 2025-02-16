CREATE TABLE users (
                    id SERIAL PRIMARY KEY,
                    username VARCHAR(255) NOT NULL UNIQUE,
                    password VARCHAR(255) NOT NULL,
                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE coins (
                    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                    amount INTEGER NOT NULL DEFAULT 0,
                    PRIMARY KEY (user_id)
);

CREATE TABLE coin_transactions (
                                id SERIAL PRIMARY KEY,
                                from_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                                to_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                                amount INTEGER NOT NULL,
                                transaction_type VARCHAR(50) NOT NULL,
                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE inventory (
                        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        item_type VARCHAR(255) NOT NULL,
                        quantity INTEGER NOT NULL DEFAULT 0,
                        PRIMARY KEY (user_id, item_type)
);
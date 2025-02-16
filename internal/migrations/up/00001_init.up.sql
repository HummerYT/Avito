CREATE TABLE users
(
    id       uuid primary key,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL,
    coins    INTEGER DEFAULT 1000 CHECK (coins >= 0)
);

CREATE TABLE transactions
(
    id           uuid PRIMARY KEY,
    from_user_id uuid REFERENCES users (id),
    to_user_id   uuid REFERENCES users (id),
    amount       INTEGER NOT NULL CHECK (amount > 0),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inventory
(
    id        uuid PRIMARY KEY,
    user_id   uuid REFERENCES users (id),
    item_type VARCHAR(255) NOT NULL,
    quantity  INTEGER DEFAULT 1 CHECK (quantity >= 0)
);
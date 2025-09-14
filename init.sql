CREATE USER main_user WITH PASSWORD 'secretpass';
GRANT ALL PRIVILEGES ON DATABASE products TO main_user;

CREATE TABLE orders (
    order_uid VARCHAR(20) PRIMARY KEY,
    track_number VARCHAR(15) NOT NULL,
    entry VARCHAR(10) NOT NULL,
    locale VARCHAR(6) NOT NULL,
    internal_signature VARCHAR(50),
    customer_id VARCHAR(25) NOT NULL,
    delivery_service VARCHAR(30) NOT NULL,
    shardkey VARCHAR(5) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard VARCHAR(5) NOT NULL
);

CREATE TABLE delivery (
    order_uid VARCHAR(20) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    "name" VARCHAR(50) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL
);

CREATE TABLE payment (
    "transaction" VARCHAR(20) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    request_id VARCHAR(20),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE items (
    order_uid VARCHAR(20) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(15) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(25) NOT NULL,
    "name" VARCHAR(50) NOT NULL,
    sale INTEGER NOT NULL,
    "size" VARCHAR(50) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(50) NOT NULL,
    status INTEGER NOT NULL,
    PRIMARY KEY (order_uid, chrt_id)
);
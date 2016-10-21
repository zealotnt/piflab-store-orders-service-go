
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TYPE order_status AS ENUM ('cart', 'processing', 'shipping', 'cancelled', 'completed');
CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
	access_token TEXT NOT NULL UNIQUE,
	status order_status DEFAULT 'cart',

	code varchar(255),
	customer_name varchar(255),
	customer_address varchar(255),
	customer_phone varchar(255),
	customer_email varchar(255),
	note varchar(255),

	created_at timestamp,
	updated_at timestamp
);
CREATE TABLE order_items (
	id SERIAL PRIMARY KEY,
	order_id SERIAL REFERENCES orders(id),
	product_id int NOT NULL,
	name varchar(255),
	price int NOT NULL DEFAULT 0,
	quantity int
);
CREATE TABLE order_status_logs (
	id SERIAL PRIMARY KEY,
	code varchar(255),
	status order_status DEFAULT 'processing',
	created_at timestamp
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE order_status_logs;
DROP TABLE order_items;
DROP TABLE orders;
DROP TYPE order_status;

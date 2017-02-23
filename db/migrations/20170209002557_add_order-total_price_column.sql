
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM order_status_logs;
DELETE FROM order_items;
DELETE FROM orders;

-- add total_price column to table orders
ALTER TABLE orders ADD total_price int;
ALTER TABLE order_status_logs ADD order_id SERIAL;
ALTER TABLE order_status_logs ADD FOREIGN KEY(order_id) REFERENCES orders(id);
ALTER TABLE order_status_logs DROP COLUMN code;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM order_status_logs;
DELETE FROM order_items;
DELETE FROM orders;

ALTER TABLE orders DROP COLUMN total_price;
ALTER TABLE order_status_logs DROP COLUMN order_id;
ALTER TABLE order_status_logs ADD code varchar(255);

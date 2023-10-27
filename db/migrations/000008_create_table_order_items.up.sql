CREATE TABLE IF NOT EXISTS order_items (
	id SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders(id),
	product_ID INT NOT NULL,
	quantity INT NOT NULL,
	price DECIMAL(50, 2)
);

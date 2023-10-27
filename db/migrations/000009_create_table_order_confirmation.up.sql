CREATE TABLE IF NOT EXISTS order_confirmation (
	id SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders(id),
	groceries_id VARCHAR(255) NOT NULL,
	confirmation_date BIGINT NOT NULL
);

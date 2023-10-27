CREATE TABLE IF NOT EXISTS orders (
	id BIGSERIAL PRIMARY KEY,
	customer_id VARCHAR(255) NOT NULL,
	total_amount DECIMAL(255,2),
	order_date BIGINT
);

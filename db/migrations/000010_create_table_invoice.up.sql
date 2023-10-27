CREATE TABLE IF NOT EXISTS invoice(
	id SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders(id),
	customer_id VARCHAR(255) NOT NULL,
	invoice_date BIGINT NOT NULL,
	total_amount DECIMAL(50,2),
	status VARCHAR(50)
);

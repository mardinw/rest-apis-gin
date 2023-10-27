CREATE TABLE IF NOT EXISTS checkouts (
	id bigserial primary key,
	customer_id varchar(255) not null,
	total_amount decimal(100, 2) not null,
	created_at bigint not null
);

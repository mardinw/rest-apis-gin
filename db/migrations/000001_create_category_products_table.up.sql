CREATE TABLE IF NOT EXISTS category_products (
	id serial primary key,
	user_id varchar(255) not null,
	name varchar(100) not null,
	description varchar(100) not null
);

CREATE INDEX idx_category_name ON category_products(name);

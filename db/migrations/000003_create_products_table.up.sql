CREATE TABLE IF NOT EXISTS products (
	id bigserial primary key,
	user_id varchar(255) not null,
	product_code varchar(255) not null,
	product_name varchar(255) not null,
	picture varchar(255),
	position varchar(255),
	quantity smallint,
	size_type_id integer,
	category_id integer,
	buy_price integer,
	mrp integer,
	defective integer,
	active bool not null,
	created bigint,
	updated bigint,
	foreign key (category_id) references category_products(id) ON UPDATE CASCADE ON DELETE CASCADE,
	foreign key (size_type_id) references size_type(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX idx_product_name ON products(product_name);

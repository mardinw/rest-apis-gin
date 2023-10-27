CREATE TABLE IF NOT EXISTS carts (
	id bigserial primary key,
	customer_id varchar(255) not null,
	product_id integer not null,
	size_type_id integer not null,
	quantity integer not null,
	comments text,
	created_at bigint not null,
	foreign key (product_id) references products(id),
	foreign key (size_type_id) references size_type(id)
);

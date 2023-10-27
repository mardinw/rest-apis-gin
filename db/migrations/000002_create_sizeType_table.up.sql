CREATE TABLE IF NOT EXISTS size_type (
	id serial primary key,
	user_id varchar(255) not null,
	name text unique not null
);

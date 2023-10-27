CREATE TABLE IF NOT EXISTS operationals (
	id bigserial primary key,
	groceries_id varchar(255) unique not null,
	day_operational text[],
	open bigint not null,
	close bigint not null,
	active bool
);

CREATE TABLE users (
	login varchar not null unique,
	password varchar not null
);

INSERT INTO users (login, password) VALUES ('admin', 'admin');


CREATE TABLE user_transaction (
	id bigserial primary key not null,
	user_id bigint not null,
	user_email varchar not null,
	amount bigint not null,
	currency varchar not null,
	created_at timestamp default now(),
	updated_at timestamp default now(),
	status varchar not null default 'НОВЫЙ',
	cancel_status bool not null default false
);

INSERT INTO user_transaction (user_id, user_email, amount, currency, created_at, updated_at) VALUES(1, 'aaa@aaa.com', 300, 'RUB', now(), now());
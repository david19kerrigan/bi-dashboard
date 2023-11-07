CREATE SCHEMA finni;

CREATE TABLE finni.patients (
	id serial primary key,
	name text NOT NULL,
	status text NOT NULL,
	address text NOT NULL
);

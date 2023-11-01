CREATE SCHEMA finni
    CREATE TABLE patients (
		id serial primary key,
		name text NOT NULL,
		status text NOT NULL,
		address text NOT NULL
	);

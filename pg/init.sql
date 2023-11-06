CREATE SCHEMA finni
    CREATE TABLE patients (
		id serial primary key,
		name text NOT NULL,
		status text NOT NULL,
		address text NOT NULL,
	);
    CREATE TABLE finni.custom (
		id serial primary key,
		names text 
	);
    CREATE TABLE finni.patients_custom (
		custom_id int references finni.custom(id),
		patient_id int references finni.patients(id),
		value text 
	);

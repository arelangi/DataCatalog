CREATE SEQUENCE datacatalog.public.metadata_id_seq;

CREATE TABLE datacatalog.public.metadata(
	dataset_id int NOT NULL DEFAULT nextval('metadata_id_seq'),
	dataset_uuid uuid DEFAULT uuid_generate_v4 (),
	dataset_name varchar NOT NULL,
	dataset_logical_name varchar,
	dataset_description varchar,
	dataset_type varchar,
	dataset_source varchar,
	dataset_share varchar,
	dataset_retention int,
	dataset_retention_justification varchar,
	dataset_arrival_frequency varchar,
	organization varchar NOT NULL,
	product varchar NOT NULL,
	team varchar NOT NULL,
	data_steward varchar NOT NULL,
	platform_name varchar NOT NULL,
	host_name varchar,
	database_name varchar,
	schema_name varchar,
	data_classiffication varchar NOT NULL,
	created_date date NOT NULL  DEFAULT CURRENT_DATE,
	last_updated_time timestamp NOT NULL  DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(dataset_id)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT uuid_generate_v1();
SELECT uuid_generate_v4();

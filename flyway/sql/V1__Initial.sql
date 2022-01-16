CREATE SCHEMA IF NOT EXISTS events;

CREATE TABLE events.processed_events (
	id uuid NOT NULL,
	processed_timestamp timestamp NOT NULL,
	event_name varchar(256) NOT NULL
);

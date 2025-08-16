-- liquibase formatted sql
-- changeset jamie:v1.1.0_create_initial_tables failOnError:true

CREATE TABLE stops
(
    id        TEXT PRIMARY KEY,
    name      TEXT,
    latitude  DECIMAL,
    longitude DECIMAL
);

CREATE TABLE routes
(
    id             TEXT,
    agency_id      TEXT,
    service_number TEXT,
    route_name     TEXT,
    route_type     TEXT,
    direction      TEXT,
    PRIMARY KEY (service_number, route_name)
);

GRANT
SELECT,
INSERT
,
UPDATE,
DELETE
ON ALL TABLES IN SCHEMA public
    TO "bus-api";
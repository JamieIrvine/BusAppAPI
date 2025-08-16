-- liquibase formatted sql
-- changeset jamie:v1.1.0_create_initial_tables failOnError:true

CREATE TABLE stops(
    id TEXT PRIMARY KEY,
    name TEXT,
    latitude DECIMAL,
    longitude DECIMAL
);

GRANT SELECT, INSERT, UPDATE, DELETE
      ON ALL TABLES IN SCHEMA public
          TO "bus-api";
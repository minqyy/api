CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(64) NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS documents (
    document_id uuid DEFAULT uuid_generate_v4() NOT NULL,
    parsed json,
    parse_error text,
    created_at timestamp DEFAULT NOW() NOT NULL,
    CONSTRAINT documents_pk PRIMARY KEY (document_id)
);

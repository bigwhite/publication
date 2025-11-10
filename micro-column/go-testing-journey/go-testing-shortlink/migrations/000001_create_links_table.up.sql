-- migrations/000001_create_links_table.up.sql
CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS short_code_idx ON links (short_code);

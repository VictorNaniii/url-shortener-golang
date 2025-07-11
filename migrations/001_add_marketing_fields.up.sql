-- Up migration: add marketing feature fields to short_urls table
ALTER TABLE short_urls
    ADD COLUMN custom_alias VARCHAR(255) UNIQUE,
    ADD COLUMN expiration TIMESTAMP NULL,
    ADD COLUMN max_clicks BIGINT NULL,
    ADD COLUMN utm_source VARCHAR(255) NULL,
    ADD COLUMN utm_medium VARCHAR(255) NULL,
    ADD COLUMN utm_campaign VARCHAR(255) NULL;


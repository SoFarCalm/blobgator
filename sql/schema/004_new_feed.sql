-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fecthed_at TIMESTAMP NULL;


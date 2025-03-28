
-- +goose Up
ALTER TABLE feeds 
ADD COLUMN last_fetched_at TIMESTAMP NULL;
-- +goose Down
ALTER TABLE feeds
DELETE COLUMN last_fetched_at;
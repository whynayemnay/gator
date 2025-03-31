-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    title TEXT NOT NULL, 
    url TEXT NOT NULL UNIQUE,
    description TEXT NULL,
    published_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_feed FOREIGN KEY (feed_id)
        REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
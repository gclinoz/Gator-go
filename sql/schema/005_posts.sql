-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR,
    url VARCHAR NOT NULL,
    description VARCHAR,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE,
    FOREIGN KEY(feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;

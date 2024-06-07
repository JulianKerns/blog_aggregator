-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(user_id,feed_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) 
);

-- +goose Down
DROP TABLE feed_follows;
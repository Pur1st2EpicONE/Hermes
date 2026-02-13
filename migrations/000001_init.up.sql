CREATE TABLE IF NOT EXISTS comments (
    id         INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    parent_id  INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    author     VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id_created_at ON comments (parent_id, created_at DESC);

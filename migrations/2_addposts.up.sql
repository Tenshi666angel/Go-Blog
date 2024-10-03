CREATE TABLE IF NOT EXISTS posts(
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    likes_count INTEGER NOT NULL DEFAULT FALSE,
    user_id INTEGER NOT NULL);

CREATE INDEX IF NOT EXISTS idx_user_id ON posts(user_id);
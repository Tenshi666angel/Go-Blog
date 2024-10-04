CREATE TABLE IF NOT EXISTS likes(
    id INTEGER PRIMARY KEY,
    post_id TEXT NOT NULL,
    username TEXT NOT NULL);

CREATE INDEX IF NOT EXISTS idx_postid ON likes(post_id);
CREATE INDEX IF NOT EXISTS idx_usernameliked ON likes(username);
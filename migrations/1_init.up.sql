CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password NOT NULL);
CREATE INDEX IF NOT EXISTS idx_username ON users(username);
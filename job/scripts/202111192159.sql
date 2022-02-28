CREATE TABLE jobs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  params TEXT NOT NULL,
  at TEXT NOT NULL,
  attempts INTEGER NOT NULL,
  max_attempts INTEGER  NOT NULL,
  locked_until TEXT,
  failed TEXT
)

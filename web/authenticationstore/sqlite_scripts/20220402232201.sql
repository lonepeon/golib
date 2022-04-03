CREATE TABLE authentication_user (
  ID TEXT PRIMARY KEY,
  username TEXT,
  password TEXT,
  salt TEXT
);

CREATE UNIQUE INDEX authentication_user_username ON authentication_user(username);

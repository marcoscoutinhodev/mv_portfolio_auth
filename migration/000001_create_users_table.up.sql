CREATE TABLE IF NOT EXISTS users (
  id varchar(36) NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email varchar(320) NOT NULL,
  password VARCHAR(50) NOT NULL
);
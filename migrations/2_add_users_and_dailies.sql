CREATE TABLE dailies (
  id serial PRIMARY KEY,
  recipes_offset int NOT NULL,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE TABLE users (
  id serial PRIMARY KEY,
  username varchar(50) NOT NULL,
  email varchar(254) NOT NULL,
  passwords varchar(40),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

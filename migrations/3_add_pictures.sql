CREATE TABLE pictures (
  id serial PRIMARY KEY,
  user_id int references users (id),
  path text NOT NULL,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

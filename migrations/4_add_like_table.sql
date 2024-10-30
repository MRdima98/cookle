CREATE TABLE likes (
  id serial PRIMARY KEY,
  user_id int references users (id),
  picture_id int references pictures (id),
  thumbs_up int DEFAULT 1,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE TABLE recipes (
  id serial PRIMARY KEY,
  name varchar(120),
  minutes int NOT NULL,
  contributor_id int NOT NULL,
  submitted date NOT NULL,
  tags varchar(5000),
  nutrition varchar(100) NOT NULL,
  n_steps int NOT NULL,
  steps varchar(20000) NOT NULL,
  description varchar(10000),
  ingredients varchar(1000) NOT NULL,
  n_ingredients int NOT NULL,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

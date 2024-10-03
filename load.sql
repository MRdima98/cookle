copy recipes (
  name,
  id,
  minutes,
  contributor_id,
  submitted,
  tags,
  nutrition,
  n_steps,
  steps,
  description,
  ingredients,
  n_ingredients
)
FROM
  '/tmp/RAW_recipes.csv' DELIMITER ',' CSV HEADER;

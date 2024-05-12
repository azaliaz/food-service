CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    mealtype VARCHAR(100),
    calories VARCHAR(100),
    protein VARCHAR(100),
    fat VARCHAR(100),
    carbohydrates VARCHAR(100),
    grams VARCHAR(100)
);
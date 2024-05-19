CREATE TABLE IF NOT EXISTS products(
    product_id serial PRIMARY KEY,
    user_id int,
    name VARCHAR(100),
    mealtype VARCHAR(100),
    calories VARCHAR(100),
    protein VARCHAR(100),
    fat VARCHAR(100),
    carbohydrates VARCHAR(100),
    grams VARCHAR(100),
    eating_date date
);

CREATE TABLE IF NOT EXISTS users(
    user_id serial primary key,
    username varchar(100),
    password varchar(100)
);

 
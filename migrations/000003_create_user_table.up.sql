CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passportNumber VARCHAR(6) UNIQUE NOT NULL,
    passportSerie VARCHAR(4) UNIQUE NOT NULL,
    surname VARCHAR(50),
    name VARCHAR(50),
    patronymic VARCHAR(50),
    address VARCHAR(100)
);



CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    age INTEGER,
    gender VARCHAR(20),
    nationality VARCHAR(50)
    );

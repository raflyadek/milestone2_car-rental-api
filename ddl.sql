CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) UNIQUE NOT NULL,
    password VARCHAR(300) NOT NULL,
    full_name VARCHAR(200) NOT NULL
    validation_code VARCHAR(300),
    validation_status BOOLEAN DEFAULT FALSE
);

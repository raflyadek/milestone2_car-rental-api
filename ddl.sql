CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) UNIQUE NOT NULL,
    password VARCHAR(300) NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    validation_code VARCHAR(300),
    deposit DECIMAL(15, 2) DEFAULT 0.0,
    validation_status BOOLEAN DEFAULT FALSE,
    role VARCHAR(20) DEFAULT 'user'
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    plat_number VARCHAR(100) NOT NULL UNIQUE,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    description TEXT NOT NULL,
    price DECIMAL(12, 2) NOT NULL CHECK (price > 0),
    availability BOOLEAN DEFAULT TRUE,
    availability_until DATE
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    car_id INT REFERENCES cars(id) ON DELETE CASCADE,
    start_date DATE DEFAULT CURRENT_DATE,
    end_date DATE NOT NULL,
    price DECIMAL(15, 2) NOT NULL CHECK (price > 0),
    status BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 10 minute
    valid_until TIMESTAMP
);

-- (if payment.status = true then add the payment data to this rental_logs)
CREATE TABLE rental_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    car_id INT REFERENCES cars(id) ON DELETE CASCADE,
    payment_id INT REFERENCES payments(id) ON DELETE CASCADE,
    total_day INT NOT NULL CHECK (total_day > 0),
    total_spent DECIMAL(15, 2) NOT NULL CHECK (total_spent > 0),
    rental_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
CREATE TABLE temperature_data (
  id SERIAL PRIMARY KEY,
  date DATE NOT NULL,
  daily_highest FLOAT NOT NULL,
  daily_lowest FLOAT NOT NULL
);

CREATE TABLE phone_numbers (
  id SERIAL PRIMARY KEY,
  phone_number VARCHAR(20) NOT NULL
);

CREATE TABLE telegram (
  id SERIAL PRIMARY KEY,
  botToken VARCHAR(45) NOT NULL 
);
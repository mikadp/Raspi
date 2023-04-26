CREATE TABLE temperature_data (
  id INT NOT NULL AUTO_INCREMENT,
  date DATE NOT NULL,
  daily_highest FLOAT NOT NULL,
  daily_lowest FLOAT NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE phone_numbers (
  id INT NOT NULL AUTO_INCREMENT,
  phone_number VARCHAR(20) NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE telegram (
  id INT NOT NULL AUTO_INCREMENT,
  botToken VARCHAR(45) NOT NULL 
);
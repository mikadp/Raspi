CREATE TABLE temperature_data (
  id INT NOT NULL AUTO_INCREMENT,
  date DATE NOT NULL,
  daily_highest FLOAT NOT NULL,
  daily_lowest FLOAT NOT NULL,
  PRIMARY KEY (id)
);

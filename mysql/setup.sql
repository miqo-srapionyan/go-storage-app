CREATE SCHEMA IF NOT EXISTS `storage_app` COLLATE = utf8mb4_0900_ai_ci;
USE `storage_app`;
CREATE TABLE IF NOT EXISTS promotions (
      id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
      uuid VARCHAR(36),
      price DOUBLE,
      expiration_date DATETIME
);
CREATE DATABASE IF NOT EXISTS user_management;
USE user_management;

-- Drop table if it exists
DROP TABLE IF EXISTS UserService;

-- Create UserService table
CREATE TABLE UserService(
    UserId INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(50) NOT NULL, 
    Email VARCHAR(50) UNIQUE NOT NULL, 
    ContactNo CHAR(8) NOT NULL, 
    Dob DATE NOT NULL, 
    Address VARCHAR(150) NOT NULL, 
    HashedPassword VARCHAR(255) NOT NULL,
    MembershipLevel ENUM('Basic', 'Premium', 'VIP') NOT NULL DEFAULT 'Basic',
    CreatedDateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert data into UserService table
INSERT INTO UserService (Name, Email, ContactNo, Dob, Address, HashedPassword, MembershipLevel, CreatedDateTime) 
VALUES 
('Mike Tan', 'xinyu@gmail.com', '99312568', '1998-10-09', '41 Woodlands Drive', '$2a$14$tUIOKXZjOjGloasiKhQU0uRUe1.b9uEYIP4/y5fiXzQMJiYDO8fe2', 'Basic', '2025-01-19 20:08:07'),
('Julie Phang', 'julie@gmail.com', '80984356', '1995-08-12', '8 Jurong West Drive', '$2a$14$tUIOKXZjOjGloasiKhQU0uRUe1.b9uEYIP4/y5fiXzQMJiYDO8fe2', 'Premium', '2024-12-20 10:38:58'),
('Joe Doe', 'joemama123@gmail.com', '98568995', '2000-11-01', '29 Clementi East Street', '$2a$14$tUIOKXZjOjGloasiKhQU0uRUe1.b9uEYIP4/y5fiXzQMJiYDO8fe2', 'VIP', '2024-11-20 23:01:09'),
('Lee Hi', 'hibye@gmail.com', '87654321', '1990-11-23', '25 Tanglin Road', '$2a$14$tUIOKXZjOjGloasiKhQU0uRUe1.b9uEYIP4/y5fiXzQMJiYDO8fe2', 'Basic', '2024-10-21 08:27:54');

-- Verify the data
SELECT * FROM UserService;

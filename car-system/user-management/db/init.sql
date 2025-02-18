CREATE DATABASE IF NOT EXISTS user_management;
USE user_management;

-- Drop table if they exist
DROP TABLE IF EXISTS Membership;
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
CreatedDateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert data into UserService table
INSERT INTO UserService (Name, Email, ContactNo, Dob, Address, HashedPassword, CreatedDateTime) 
VALUES 
('Mike Tan', 'miketan98@gmail.com', '99312568', '1998-10-09', '41 Woodlands Drive', '$2a$14$PZpt2VH60evaYc8LUAMth.IXOcunXnp.2t/DMJJI1RausohfuFra2', '2025-01-19 20:08:07'),
('Julie Phang', 'julie@gmail.com', '80984356', '1995-08-12', '8 Jurong West Drive', '$2a$14$fS1.LWaZnVt1QyUt4WJjMu2ZhaUEeahZFIw8sC28CiPeszZASBlG.', '2024-12-20 10:38:58'),
('Joe Doe', 'joemama123@gmail.com', '98568995', '2000-11-01', '29 Clementi East Street', '$2a$14$rJaPEsOCBgQIAtaw4b33A.apLCm.msGndS2x0RYB4Ql7bVMe/LDiu', '2024-11-20 23:01:09'),
('Lee Hi', 'hibye@gmail.com', '87654321', '1990-11-23', '25 Tanglin Road', '$2a$14$jckIwIOXPhLcqd.DNuWm3.4HfH7FlaB6yqPPLl8UOKswZ6Tdmq04C', '2024-10-21 08:27:54');

-- Create Membership table
CREATE TABLE Membership(
MembershipId INT AUTO_INCREMENT PRIMARY KEY,
UserId INT NOT NULL,
MembershipTier ENUM('Basic', 'Premium', 'Vip') NOT NULL,
HourlyRate DECIMAL(10,2) NOT NULL,
MemberDiscount DECIMAL(5,2) NOT NULL,
PriorityLevel TINYINT NOT NULL CHECK (PriorityLevel BETWEEN 0 AND 2),
TotalCostPerMonth DECIMAL(10,2),
MembershipExpiryDate DATE NOT NULL,
EligibleForUpgradeNextMonth BOOLEAN NOT NULL DEFAULT FALSE,

FOREIGN KEY (UserId) REFERENCES UserService(UserId) ON DELETE CASCADE
);

-- Insert data into Membership table
INSERT INTO Membership (UserId, MembershipTier, HourlyRate, MemberDiscount, PriorityLevel, TotalCostPerMonth, MembershipExpiryDate, EligibleForUpgradeNextMonth) 
VALUES
(1, 'Premium', 10.00, 5.00, 1, 54.00, '2025-03-20', FALSE), 
(2, 'VIP', 5.00, 10.00, 2, 22.50, '2025-04-19', FALSE), 
(3, 'Premium', 10.00, 5.00, 1, 36.00, '2025-05-19', FALSE), 
(4, 'Basic', 15.00, 0.00, 0, 108.00, '2025-06-18', FALSE); 
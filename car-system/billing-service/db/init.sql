CREATE DATABASE IF NOT EXISTS billing_sevice;
USE billing_service;

-- Drop table if they exist
DROP TABLE IF EXISTS PaymentTransactions;
DROP TABLE IF EXISTS Promotion;
DROP TABLE IF EXISTS Billing;

-- Create Billing table
CREATE TABLE Billing (
BillingId INT AUTO_INCREMENT PRIMARY KEY,
ReservationId INT NOT NULL, 
UserId INT NOT NULL,
BillingTotal DECIMAL(10,2) NOT NULL,
MembershipDiscount DECIMAL(10,2) DEFAULT 0,
PromoDiscount DECIMAL(10,2) DEFAULT 0,
BillingFinal DECIMAL(10,2) NOT NULL,
PaymentStatus ENUM('Pending', 'Paid', 'Overdue') NOT NULL DEFAULT 'Pending', 
BillingDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
FOREIGN KEY (ReservationId) REFERENCES vehicle_service.Reservation(ReservationId) ON DELETE CASCADE,
FOREIGN KEY (UserId) REFERENCES user_management.UserService(UserId) ON DELETE CASCADE
);

-- Insert data into Billing table
INSERT INTO Billing (UserId, ReservationId, BillingDate, BillingTotal, MembershipDiscount, PromoDiscount, BillingFinal, PaymentStatus) 
VALUES
(1, 1, '2024-11-20 14:00:00', 19.00, 5.00, 0.00, 14.00, 'Paid'), -- Premium, 5% discount
(1, 2, '2024-12-01 10:00:00', 54.00, 5.00, 5.00, 44.00, 'Paid'), -- Premium, 5% + 5% promo
(3, 3, '2024-12-05 14:00:00', 36.00, 5.00, 0.00, 31.00, 'Paid'), -- Premium, 5% discount
(2, 4, '2024-12-15 14:30:00', 22.50, 10.00, 0.00, 12.50, 'Paid'), -- VIP, 10% discount
(3, 5, '2024-12-19 20:00:00', 38.00, 5.00, 0.00, 33.00, 'Pending'), -- Premium, 5% discount
(4, 6, '2024-12-24 09:00:00', 108.00, 0.00, 20.00, 88.00, 'Overdue'); -- Basic, no discount, but $20 promo

-- Create Promotion table    
CREATE TABLE Promotion(
PromoId INT AUTO_INCREMENT PRIMARY KEY,
PromoType ENUM('Flat Discount', 'Percentage Discount', 'Seasonal Offer') NOT NULL,
PromoDiscountPercentage DECIMAL(5,2) NOT NULL,
PromoStartDate DATE,
PromoEndDate DATE
);

-- Insert data into Promotion table
INSERT INTO Promotion (PromoType, PromoDiscountPercentage, PromoStartDate, PromoEndDate) 
VALUES
('Flat Discount', 10.00, '2024-12-15', '2024-12-31'),
('Percentage Discount', 5.00, '2024-12-01', '2024-12-10'),
('Seasonal Offer', 20.00, '2024-12-20', '2024-12-25');

-- Create PaymentTransactions table
CREATE TABLE PaymentTransactions (
TransactionId INT AUTO_INCREMENT PRIMARY KEY,
BillingId INT NOT NULL,
UserId INT NOT NULL,
PaymentMethod ENUM('Credit Card', 'Debit Card', 'PayNow', 'BankTransfer') NOT NULL,
PaymentStatus ENUM('Pending', 'Completed', 'Failed') DEFAULT 'Pending',
TransactionAmount DECIMAL(10,2) NOT NULL, 
TransactionDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
FOREIGN KEY (BillingId) REFERENCES billing_service.Billing(BillingId) ON DELETE CASCADE,
FOREIGN KEY (UserId) REFERENCES user_management.UserService(UserId) ON DELETE CASCADE
);

-- Insert data into PaymentTransactions table
INSERT INTO PaymentTransactions (BillingId, UserId, PaymentMethod, PaymentStatus, TransactionAmount, TransactionDate) 
VALUES
(1, 1, 'Credit Card', 'Completed', 14.00, '2024-11-21 14:00:00'),
(2, 1, 'Debit Card', 'Completed', 44.00, '2024-12-01 20:00:00'),
(3, 3, 'Credit Card', 'Completed', 31.00, '2024-12-07 16:00:00'),
(4, 2, 'PayNow', 'Completed', 12.50, '2024-12-17 14:30:00'),
(5, 3, 'BankTransfer', 'Pending', 33.00, '2024-12-19 22:00:00'),
(6, 4, 'Credit Card', 'Failed', 88.00, '2024-12-24 11:00:00');


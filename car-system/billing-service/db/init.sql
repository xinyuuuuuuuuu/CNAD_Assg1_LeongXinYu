-- Create the billing_service database
CREATE DATABASE IF NOT EXISTS billing_service;
USE billing_service;

-- Drop tables if they exist
DROP TABLE IF EXISTS Promotion;
DROP TABLE IF EXISTS Billing;

-- Create the Billing table
CREATE TABLE Billing (
    BillingId INT AUTO_INCREMENT PRIMARY KEY,
    UserId INT NOT NULL,
    ReservationId INT NOT NULL, 
    BillingTotal DECIMAL(10,2) NOT NULL,
    PromoDiscount DECIMAL(10,2) DEFAULT 0,
    BillingFinal DECIMAL(10,2) NOT NULL,
    PaymentStatus ENUM('Pending', 'Paid', 'Overdue') NOT NULL DEFAULT 'Pending', 
    PaymentMethod ENUM('Card', 'PayNow', 'Nil') NOT NULL,
    BillingDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the Promotion table
CREATE TABLE Promotion(
    PromoId INT AUTO_INCREMENT PRIMARY KEY,
    PromoType ENUM('Flat Discount', 'Percentage Discount', 'Seasonal Offer') NOT NULL,
    PromoDiscountPercentage DECIMAL(5,2) NOT NULL,
    PromoStartDate DATE,
    PromoEndDate DATE
);

INSERT INTO Promotion (PromoType, PromoDiscountPercentage, PromoStartDate, PromoEndDate) 
VALUES
('Flat Discount', 10.00, '2024-12-15', '2024-12-31'),
('Percentage Discount', 5.00, '2024-12-01', '2024-12-10'),
('Seasonal Offer', 20.00, '2024-12-20', '2024-12-25');

-- Insert Billing records
INSERT INTO Billing (UserId, ReservationId, BillingDate, BillingTotal, PromoDiscount, BillingFinal, PaymentStatus, PaymentMethod)
VALUES
(1, 1, '2024-12-10 08:43:05', 40.00, 0.00, 40.00, 'Paid', 'Card'),
(1, 2, '2024-12-06 07:43:05', 60.00, 3.00, 57.00, 'Paid', 'Card'),
(3, 3, '2025-02-18 10:43:05', 40.00, 0.00, 40.00, 'Pending', 'Nil'),
(2, 4, '2025-02-18 07:43:05', 30.00, 3.00, 27.00, 'Pending', 'PayNow'),
(3, 5, '2025-03-15 10:43:05', 40.00, 0.00, 40.00, 'Pending', 'Card'),
(4, 6, '2025-03-06 09:43:05', 135.00, 27.00, 108.00, 'Pending', 'PayNow');




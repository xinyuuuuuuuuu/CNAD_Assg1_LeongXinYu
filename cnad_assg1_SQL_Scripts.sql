CREATE DATABASE cnad_assg1
USE cnad_assg1;

CREATE TABLE UserService(
UserId CHAR(5) NOT NULL PRIMARY KEY, 
Name VARCHAR(50) NOT NULL, 
Email VARCHAR(50) NOT NULL, 
ContactNo CHAR(8) NOT NULL, 
Dob DATETIME NOT NULL, 
Address VARCHAR(150) NOT NULL, 
Password VARCHAR(255) NOT NULL,
CreatedDateTime DATETIME NOT NULL
);

CREATE TABLE Vehicle(
VehicleId CHAR(5) NOT NULL PRIMARY KEY,
VehicleMake VARCHAR(50) NOT NULL,
VehicleModel VARCHAR(50) NOT NULL ,
VehicleType VARCHAR(30) NOT NULL,
LicensePlate VARCHAR(15) NOT NULL,
VehicleStatus CHAR(2) NOT NULL CHECK (VehicleStatus IN ('R', 'A', 'NA')),
VehicleLocation VARCHAR(255) NOT NULL,
VehicleChargeLevel TINYINT UNSIGNED NOT NULL,   -- UNSIGNED ensure values are non negative
VehicleCleanliness VARCHAR(100) NOT NULL
);

CREATE TABLE Promotion(
PromoId CHAR(5) NOT NULL PRIMARY KEY,
PromoType CHAR(1) NOT NULL CHECK (PromoType IN ('F', 'O', 'S')),
PromoDiscount DECIMAL(5,2) NOT NULL,
PromoStartDate DATETIME NOT NULL,
PromoEndDate DATETIME NOT NULL
);

CREATE TABLE Membership(
MembershipId CHAR(5) NOT NULL PRIMARY KEY,
UserId CHAR(5) NOT NULL,
MembershipTier CHAR(10) NOT NULL CHECK (MembershipTier IN ('Basic', 'Premium', 'Vip')),
HourlyRate DECIMAL(10,2) NOT NULL,
MemberDiscount DECIMAL(5,2) NOT NULL,
PriorityLevel CHAR(1) NOT NULL CHECK (PriorityLevel IN ('0', '1', '2')),
TotalCostPerMonth DECIMAL(10,2) NULL,
MembershipExpiryDate DATETIME NOT NULL,
EligibleForUpgradeNextMonth CHAR(1) NOT NULL CHECK (EligibleForUpgradeNextMonth IN ('T', 'F')),

FOREIGN KEY (UserId)
 REFERENCES UserService(UserId)
);

CREATE TABLE Reservation(
ReservationId CHAR(5) NOT NULL PRIMARY KEY,
UserId CHAR(5) NOT NULL,
VehicleId CHAR(5) NOT NULL,
ReserveStatus CHAR(5) NOT NULL CHECK (ReserveStatus IN ('Pend', 'Conf', 'Canc')),
ReserveStartDate DATETIME NOT NULL,
ReserveEndDate DATETIME NOT NULL,
EstimatedTotalCost DECIMAL(10,2) NULL,
CreatedDate DATETIME NOT NULL,
ModifiedDate DATETIME NULL,

FOREIGN KEY (UserId) 
 REFERENCES UserService(UserId), 
FOREIGN KEY (VehicleId) 
 REFERENCES Vehicle(VehicleId) 
);

CREATE TABLE Billing(
BillingId CHAR(5) NOT NULL PRIMARY KEY,
UserId CHAR(5) NOT NULL,
ReservationId CHAR(5) NOT NULL, 
BillingDate DATETIME NOT NULL,
BillingTotal DECIMAL (10,2) NOT NULL,
PaymentMethod CHAR(2) NOT NULL CHECK (PaymentMethod IN ('CC', 'DB', 'UP')), 
PaymentStatus CHAR(2) NOT NULL CHECK (PaymentStatus IN ('S', 'NS', 'UP')),

FOREIGN KEY (UserId) 
 REFERENCES UserService(UserId),
FOREIGN KEY (ReservationId) 
 REFERENCES Reservation(ReservationId)
);

INSERT INTO UserService (UserId, Name, Email, Password, ContactNo, Dob, Address, CreatedDateTime) 
VALUES ("U0001", "Mike Tan", "miketan98@gmail.com", "$2a$1$9tdXvvEBXIrVJzN1MW/p3.8NHnVRoHcCbTaoRcXOSrBWNdbi8swqu",  "99312568", "1998-10-09", "41 Woodlands Drive", "2024-05-07 20:08:07"),
("U0002", "Julie Phang", "julie@gmail.com", "$2a$14$BmviTGH5SYRWA6rArtbVlON1F8Fr4a7spAe4O2HdsLdJyJN3uYcUi", "80984356", "1995-08-12", "8 Jurong West Drive", "2024-09-29 10:38:58"),
("U0003", "Joe Doe", "joemama123@gmail.com", "$2a$14$znzL0QpLh055BjgAwtvbCuO7.smxOQ/6gMlxC4f3LCEPzUy8X0W2O", "98568995", "2000-11-01", "29 Clementi East Street", "2024-11-09 23:01:09"),
("U0004", "Lee Hi", "hibye@gmail.com", "$2a$14$tsaEz9/VeQruHUFZXOaD/OBMmoPQK4ZEaA1HU22aMVboHeAu9xq4G", "87654321", "1990-11-23", "8 Jurong West Drive", "2024-07-13 08:27:54");

INSERT INTO Vehicle (VehicleId, VehicleMake, VehicleModel, VehicleType, LicensePlate, VehicleStatus, VehicleLocation, VehicleChargeLevel, VehicleCleanliness) 
VALUES('V0001', 'Tesla', 'Model S', 'Sedan', 'ABC1234', 'A', 'Downtown Parking Lot 1', 85, 'Clean'),
('V0002', 'Nissan', 'Leaf', 'Hatchback', 'XYZ5678', 'R', 'Uptown Garage Level 3', 70, 'Moderate'),
('V0003', 'BMW', 'i3', 'Sedan', 'DEF9012', 'NA', 'Central Mall Basement 2', 50, 'Dirty'),
('V0004', 'Chevrolet', 'Bolt EV', 'Compact', 'GHI3456', 'A', 'Airport Terminal 2 Lot A', 95, 'Clean'),
('V0005', 'Hyundai', 'Kona Electric', 'SUV', 'JKL7890', 'R', 'Suburban Plaza Lot C', 60, 'Moderate');

INSERT INTO Promotion (PromoId, PromoType, PromoDiscount, PromoStartDate, PromoEndDate) 
VALUES('P0001', 'F', 10.00, '2024-12-15 00:00:00', '2024-12-31 23:59:59'),
('P0002', 'O', 5.00, '2024-12-01 00:00:00', '2024-12-10 23:59:59'),
('P0003', 'S', 20.00, '2024-12-20 00:00:00', '2024-12-25 23:59:59');


INSERT INTO Membership (MembershipId, UserId, MembershipTier, HourlyRate, MemberDiscount, PriorityLevel, TotalCostPerMonth, MembershipExpiryDate, EligibleForUpgradeNextMonth ) 
VALUES
('M0001', 'U0001', 'Premium', 10.00, 5.00, 1, 54.00, '2025-02-01 10:00:00', 'F'), 
('M0002', 'U0002', 'VIP', 5.00, 10.00, 2, 22.50, '2025-02-15 14:30:00', 'F'), 
('M0003', 'U0003', 'Premium', 10.00, 5.00, 1, 36.00, '2025-02-05 14:00:00', 'F'), 
('M0004', 'U0004', 'Basic', 15.00, 0.00, 0, 108.00, '2025-02-24 09:00:00', 'F'); 


INSERT INTO Reservation (ReservationId, UserId, VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost, CreatedDate, ModifiedDate) 
VALUES
('R0001', 'U0001', 'V0005', 'Canc', '2024-11-21 14:00:00', '2024-11-21 16:00:00', 20.00, '2024-11-20 14:00:00', '2024-11-20 15:00:00'),  
('R0002', 'U0001', 'V0001', 'Conf', '2024-12-01 12:00:00', '2024-12-01 18:00:00', 60.00, '2024-12-01 10:00:00', '2024-12-01 11:00:00'), 
('R0003', 'U0003', 'V0002', 'Conf', '2024-12-06 14:00:00', '2024-12-06 18:00:00', 40.00, '2024-12-05 14:00:00', '2024-12-05 15:00:00'), 
('R0004', 'U0002', 'V0002', 'Conf', '2024-12-15 16:00:00', '2024-12-15 21:00:00', 25.00, '2024-12-15 14:30:00', '2024-12-15 15:30:00'), 
('R0005', 'U0003', 'V0003', 'Pend', '2024-12-20 08:00:00', '2024-12-20 12:00:00', 40.00, '2024-12-19 20:00:00', NULL), 
('R0006', 'U0004', 'V0004', 'Conf', '2024-12-25 09:00:00', '2024-12-25 18:00:00', 135.00, '2024-12-24 09:00:00', '2024-12-24 10:00:00'); 

INSERT INTO Billing (BillingId, UserId, ReservationId, BillingDate, BillingTotal, PaymentMethod, PaymentStatus) 
VALUES
('B0001', 'U0001', 'R0001', '2024-11-20 14:00:00', 19.00, 'DB', 'NS' ), -- prem, 5
('B0002', 'U0001', 'R0002', '2024-12-01 10:00:00', 54.00, 'CC', 'S'), -- prem, 5% + 5% = 10
('B0003', 'U0003', 'R0003', '2024-12-05 14:00:00', 36.00, 'DB', 'S'), -- prem, 5+5 = 10
('B0004', 'U0002', 'R0004', '2024-12-15 14:30:00', 22.50, 'DB', 'S'), -- vip, 10%
('B0005', 'U0003', 'R0005', '2024-12-19 20:00:00', 38.00, 'UP', 'UP'), -- prem, 5
('B0006', 'U0004', 'R0006', '2024-12-24 09:00:00', 108.00, 'CC', 'S'); -- basic, 0%+20 = 20
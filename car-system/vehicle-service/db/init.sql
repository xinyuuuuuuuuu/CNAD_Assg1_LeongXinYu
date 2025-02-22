CREATE DATABASE IF NOT EXISTS vehicle_service;
USE vehicle_service;

-- Drop table if they exist
DROP TABLE IF EXISTS Reservation;
DROP TABLE IF EXISTS Vehicle;

CREATE TABLE Vehicle(
VehicleId INT AUTO_INCREMENT PRIMARY KEY,
VehicleMake VARCHAR(50) NOT NULL,
VehicleModel VARCHAR(50) NOT NULL ,
VehicleType VARCHAR(30) NOT NULL,
LicensePlate VARCHAR(15) UNIQUE NOT NULL,
VehicleStatus ENUM('Reserved', 'Available', 'Not Available') NOT NULL,
VehicleLocation VARCHAR(255) NOT NULL,
VehicleChargeLevel TINYINT UNSIGNED NOT NULL CHECK (VehicleChargeLevel BETWEEN 0 AND 100),   -- UNSIGNED ensure values are non negative
VehicleCleanliness ENUM('Clean', 'Moderate', 'Dirty') NOT NULL
);

INSERT INTO Vehicle (VehicleMake, VehicleModel, VehicleType, LicensePlate, VehicleStatus, VehicleLocation, VehicleChargeLevel, VehicleCleanliness) 
VALUES
('Tesla', 'Model S', 'Sedan', 'ABC1234', 'Available', 'Downtown Parking Lot 1', 85, 'Clean'),
('Nissan', 'Leaf', 'Hatchback', 'XYZ5678', 'Reserved', 'Uptown Garage Level 3', 70, 'Moderate'),
('BMW', 'i3', 'Sedan', 'DEF9012', 'Not Available', 'Central Mall Basement 2', 50, 'Dirty'),
('Chevrolet', 'Bolt EV', 'Compact', 'GHI3456', 'Available', 'Airport Terminal 2 Lot A', 95, 'Clean'),
('Hyundai', 'Kona Electric', 'SUV', 'JKL7890', 'Reserved', 'Suburban Plaza Lot C', 60, 'Moderate'),
('Ford', 'Mustang Mach-E', 'SUV', 'MNO4321', 'Available', 'Downtown Parking Lot 2', 90, 'Clean'),
('Volkswagen', 'ID.4', 'SUV', 'PQR5678', 'Available', 'City Center Garage Level 1', 80, 'Clean'),
('Audi', 'e-tron', 'SUV', 'STU6789', 'Available', 'Tech Park Basement 3', 75, 'Moderate'),
('Porsche', 'Taycan', 'Sedan', 'VWX7890', 'Available', 'Luxury Plaza VIP Parking', 65, 'Clean'),
('Mercedes-Benz', 'EQC', 'SUV', 'YZA9012', 'Available', 'Business District Parking A', 88, 'Clean');


CREATE TABLE Reservation(
ReservationId INT AUTO_INCREMENT PRIMARY KEY,
UserId INT NOT NULL,
VehicleId INT NOT NULL,
ReserveStatus ENUM('Pending', 'Active', 'Completed', 'Cancelled') NOT NULL,
ReserveStartDate DATETIME NOT NULL,
ReserveEndDate DATETIME NOT NULL,
EstimatedTotalCost DECIMAL(10,2),
CreatedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ModifiedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

FOREIGN KEY (VehicleId) REFERENCES Vehicle(VehicleId) ON DELETE CASCADE
);

INSERT INTO Reservation (UserId, VehicleId, ReserveStatus, ReserveStartDate, ReserveEndDate, EstimatedTotalCost, CreatedDate, ModifiedDate) 
VALUES
(1, 5, 'Completed', '2024-12-10 05:43:05', '2024-12-10 08:43:05', 150.00, '2024-12-09 05:43:05', '2024-12-10 08:43:05'),  
(1, 1, 'Completed', '2024-12-06 05:43:05', '2024-12-06 07:43:05', 100.00, '2024-12-05 05:43:05', '2024-12-06 07:43:05'), 
(3, 2, 'Completed', '2025-02-17 22:43:05', '2025-02-18 10:43:05', 360.00, '2025-02-16 22:43:05', '2025-02-18 10:43:05'), 
(2, 2, 'Completed', '2025-02-17 22:43:05', '2025-02-18 07:43:05', 360.00, '2025-02-16 22:43:05', '2025-02-18 07:43:05'), 
(3, 3, 'Pending', '2025-03-15 05:43:05', '2025-03-15 10:43:05', 150.00, '2025-03-14 05:43:05', NULL), 
(4, 4, 'Pending', '2025-03-06 05:43:05', '2025-03-06 09:43:05', 200.00, '2025-03-05 05:43:05', '2025-03-06 09:43:05');


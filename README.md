# CNAD_Assg1_LeongXinYu

1. Design considerations
<br><b>1.1 Loose Coupling</b></br>
   The services are independent, it interacts through APIs, minimizing dependencies on one another.

<br><b>1.2 Security</b></br>
  Bcrypt is used for password.
  JWT token is used for user-level access. Each user can only access their own information.

<br><b>1.3 Business function</b></br>
  Each microservice handles a specific business function for example, user management, vehicle service and billing service.
  User Management Service – Handles user authentication and account management.
  Vehicle Service – Manages vehicle listings and reservations.
  Billing Service – Handles payments and billing transactions.

<br><b>2. Architecture diagram</b></br>
   ![CNAD-Assg1-XinYu-ArchitectureDiagram](https://github.com/user-attachments/assets/ca7527bb-20cc-4e9e-9444-95a464fc00ec)

<br><b>3. Setting up and running the microservices</b></br>
<br><b>3.1 Environment variables</b></br>
   DB_USER=electric-car-sharing-admin
   DB_PASSWORD=Password123
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=user_management
   JWT_SECRET=bWltMnpqNGw1ejV6cG45NWJrNGt6OTdybnpvMGhqZnQ0MnRyM2w4NGF1Z3pjMzZtdG1zcnY1OHpjczNkYXpldnpqZjk3anUwMDZoZGR3Z210a3ZxczI4dmZ1MmtsZ25uZTVxZDBjMGZyeTRkaDF2b2x4aHd1d2w0eWk2ZWxxMmg=
   ![image](https://github.com/user-attachments/assets/580ae8ae-a8ef-4f94-b0e5-5190d2e6da9f)
   Insert the .env file into the root folder for user-management, vehicle-servcie, billing-service.

<br><b>3.2 Set up Database</b></br>
<br><b>3.2.1 Create account to access the database</b></br>
CREATE USER 'electric-car-sharing-admin'@'%' IDENTIFIED BY 'Password123';
GRANT ALL PRIVILEGES ON user_management.* TO 'electric-car-sharing-admin'@'%';
GRANT ALL PRIVILEGES ON vehicle_service.* TO 'electric-car-sharing-admin'@'%';
GRANT ALL PRIVILEGES ON billing_service.* TO 'electric-car-sharing-admin'@'%';
FLUSH PRIVILEGES;

<br><b>3.2.2 Database</b></br>
Use the provided init.sql in each microservice db folder to create and populate each microservice's database

<br><b>3.2.3 Run application</b></br>
cd car-system
cd user-management
go run .\main.go
cd vehicle-service
go run .\main.go
cd billing-service
go run .\main.go

<br><b>4. Login Account Password</b></br>
   Account Password: Password


   

# CNAD_Assg1_LeongXinYu

4.1.3.1. Design consideration of your microservices 
Microservices are being used in this system.
I have:
User Service, Vehicle Service, Billing Service, Reservation Service, Promotion Service.

The database used for this project is MySQL, a relational database. The tables are separated, and the primary and foreign keys are defined clearly in the sql scripts file located in this repo. 

For security, I have used bcrypt to encrpt the password to ensure that in any case that the data may leak, unauthorized user do not have access to the password of the users of this system. 

Validation has been done for each functions, however some are not working. 

From the terminal, you can run the main.go file and a menu console will be shown. User can choose to sign up for a new account or login to their current account, after that they can access a menu console specially for members of this system.


4.1.3.2. Architecture diagram 
![alt text](image.png)

4.1.3.3. Instructions for setting up and running your microservices
To be done on the terminal:
cd services,
cd userService,
go run main.go,
User Console menu will show up and user can select what option they would like.
Only after opt 2 (Login to account), user can then access 
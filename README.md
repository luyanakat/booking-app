
# Booking App

A booking-hotel website to making reservations.

Temporarily deploy on: [Ngrok](https://708e-42-117-149-19.ap.ngrok.io/) (I'm poor)

Admin Account: admin@admin.com

Password: admin

![web](2image.png)

### Build and Run Project

#### Create database

Firstly, because the program uses Postgresql to store user data, we need to create and run the database

- Next, we need to create a database named “bookings" with the following comand:

  ```sql
  CREATE DATABASE bookings;
  ```

- Install soda-cli to run migration:

  `go install github.com/gobuffalo/pop/v6/soda@latest`

- Run `soda migrate` with terminal in project-folder

- All the needed tables & constrains now created

![db](image.png)

#### Run locally

- Edit `.env`

- Simple run `run.bat` on windows or `run.sh` on linux, make sure you executable `run.sh` on linux by using `chmod +x`

App will run in  `localhost:3000`

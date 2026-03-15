# Go Users API

REST API built with Go, PostgreSQL, Redis and Docker.

This project is a backend API that manages users and demonstrates a modern backend architecture using containers, caching and a relational database.

---

## Tech Stack

* Go
* Chi Router
* PostgreSQL
* Redis
* Docker
* Docker Compose

---

## Architecture

Client → Go API → Redis Cache → PostgreSQL Database

---

## Run the Project

Clone the repository:

git clone https://github.com/Omen77796/go-users-api.git

Enter the project directory:

cd go-users-api

Start all services with Docker:

docker compose up --build

The API will run at:

http://localhost:8080

---

## API Endpoints

Health check

GET /health

Get all users

GET /users

Get user by ID

GET /users/{id}

Create new user

POST /users

---

## Example Request

POST /users

{
"name": "Joseph",
"email": "[joseph@email.com](mailto:joseph@email.com)"
}

---

## Project Structure

cmd/api
Application entry point

internal/handlers
HTTP handlers

internal/models
Data models

internal/middleware
Middleware logic

init/init.sql
Database initialization script

---

## Author

Joseph Rodriguez

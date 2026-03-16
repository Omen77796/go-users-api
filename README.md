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

## mermaid
graph TD

Client --> API[Go API - Chi Router]

API --> Redis[(Redis Cache)]

API --> Postgres[(PostgreSQL Database)]

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

## API Usage Examples

### Health Check

Request

curl http://localhost:8080/health

Response

{
"status": "ok"
}

---

### Get All Users

Request

curl http://localhost:8080/users

Example Response

[
{
"id": 1,
"name": "Omen77796",
"email": "[Omen77796@email.com](mailto:Omen77796@email.com)"
}
]

---

### Get User by ID

Request

curl http://localhost:8080/users/1

Example Response

{
"id": 1,
"name": "Omen77796",
"email": "[Omen77796@email.com](mailto:Omen77796@email.com)"
}

---

### Create User

Request

curl -X POST http://localhost:8080/users 
-H "Content-Type: application/json" 
-d '{
"name": "Omen77796",
"email": "[Omen77796@email.com](mailto:Omen77796@email.com)"
}'

Example Response

{
"id": 2,
"name": "Omen77796",
"email": "[Omen77796@email.com](mailto:Omen77796@email.com)"
}


---

## Example Request

POST /users

{
"name": "Omen77796",
"email": "[Omen77796@email.com](mailto:Omen77796@email.com)"
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

## Maintainer

Omen77796 

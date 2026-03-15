package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"go-users-api/internal/handlers"
	"go-users-api/internal/middleware"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

func main() {

	connStr := "host=postgres user=postgres password=postgres2026 dbname=usersdb sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}

		log.Println("Esperando PostgreSQL...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL:", err)
	}

	log.Println("Conectado a PostgreSQL 🚀")

	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("No se pudo conectar a Redis:", err)
	}

	log.Println("Conectado a Redis 🚀")

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", handlers.HealthHandler)
	r.Get("/users/{id}", handlers.GetUserByIDHandler(db))
	r.Method("GET", "/users", handlers.UsersHandler(db, rdb))
	r.Method("POST", "/users", handlers.UsersHandler(db, rdb))

	log.Println("Servidor corriendo en puerto 8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", r))

	log.Println("Servidor corriendo en puerto 8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

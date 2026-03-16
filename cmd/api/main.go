package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"go-users-api/internal/handlers"
	"go-users-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

func main() {

	connStr := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"

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
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
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

	port := os.Getenv("SERVER_PORT")

	log.Println("Servidor corriendo en puerto", port, "🚀")
	log.Fatal(http.ListenAndServe(":"+port, r))

}

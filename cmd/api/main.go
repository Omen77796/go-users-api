package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/omen77796/go-users-api/internal/handlers"
	"github.com/omen77796/go-users-api/internal/middleware"
	"github.com/omen77796/go-users-api/internal/repository"
	"github.com/omen77796/go-users-api/internal/services"

	_ "github.com/omen77796/go-users-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

var ctx = context.Background()

func main() {

	// ================================
	// 🔹 PostgreSQL
	// ================================
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL no está definida")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		log.Println("Esperando PostgreSQL...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL:", err)
	}

	log.Println("Conectado a PostgreSQL 🚀")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ================================
	// 🔹 Redis
	// ================================
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	for i := 0; i < 10; i++ {
		if err = rdb.Ping(ctx).Err(); err == nil {
			break
		}
		log.Println("Esperando Redis...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("No se pudo conectar a Redis:", err)
	}

	log.Println("Conectado a Redis 🚀")

	// ================================
	// 🔥 AQUI VA TU NUEVA ARQUITECTURA
	// ================================
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, rdb)
	userHandler := handlers.NewUserHandler(userService)

	// ================================
	// 🔹 Router
	// ================================
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/health", handlers.HealthHandler)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// 🔥 NUEVAS RUTAS (LIMPIAS)
	r.Get("/users", userHandler.GetUsers)
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	// ================================
	// 🔹 Server
	// ================================
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor corriendo en puerto", port, "🚀")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/omen77796/go-users-api/internal/config"
	"github.com/omen77796/go-users-api/internal/handlers"
	"github.com/omen77796/go-users-api/internal/logger"
	"github.com/omen77796/go-users-api/internal/middleware"
	"github.com/omen77796/go-users-api/internal/repository"
	"github.com/omen77796/go-users-api/internal/services"

	_ "github.com/omen77796/go-users-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	// 🔥 CARGAR CONFIG PRIMERO
	cfg := config.Load()

	// 🔥 LOGGER
	logger.Init()
	defer logger.Log.Sync()

	// ================================
	// 🔹 PostgreSQL
	// ================================
	dbURL := cfg.DBUrl
	if dbURL == "" {
		logger.Log.Fatal("DATABASE_URL no está definida")
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		logger.Log.Fatal("failed to open database", zap.Error(err))
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		logger.Log.Warn("waiting for PostgreSQL...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Log.Fatal("PostgreSQL connection failed", zap.Error(err))
	}

	logger.Log.Info("PostgreSQL connected")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ================================
	// 🔹 Redis
	// ================================
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	for i := 0; i < 10; i++ {
		if err = rdb.Ping(context.Background()).Err(); err == nil {
			break
		}
		logger.Log.Warn("waiting for Redis")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Log.Fatal("Redis connection failed", zap.Error(err))
	}

	logger.Log.Info("Redis connected")

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

	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery)
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
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// correr servidor en goroutine
	go func() {
		logger.Log.Info("server started", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("fatal error", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logger.Log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("error shutting down server", zap.Error(err))
	}

	logger.Log.Info("Server stopped")
}

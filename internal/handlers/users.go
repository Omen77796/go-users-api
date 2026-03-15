package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-users-api/internal/models"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func UsersHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		switch r.Method {

		case http.MethodGet:
			handleGetUsers(w, db, rdb)

		case http.MethodPost:
			createUserHandler(w, r, db, rdb)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		u.Name, u.Email,
	).Scan(&u.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rdb.Del(ctx, "users")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func handleGetUsers(w http.ResponseWriter, db *sql.DB, rdb *redis.Client) {

	cached, err := rdb.Get(ctx, "users").Result()

	if err == nil {
		w.Write([]byte(cached))
		log.Println("Datos obtenidos desde Redis ⚡")
		return
	}

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var u models.User

		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, u)
	}

	func GetUserByIDHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		var user models.User

		err := db.QueryRow(
			"SELECT id, name, email FROM users WHERE id = $1",
			id,
		).Scan(&user.ID, &user.Name, &user.Email)

		if err != nil {

			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rdb.Set(ctx, "users", jsonData, 30*time.Second)

	log.Println("Datos obtenidos desde PostgreSQL 🐘")

	w.Write(jsonData)
}

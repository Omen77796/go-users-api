package repository

import (
	"database/sql"

	"github.com/omen77796/go-users-api/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// CREATE
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		user.Name, user.Email,
	).Scan(&user.ID)
}

// GET ALL
func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GET BY ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User

	err := r.DB.QueryRow(
		"SELECT id, name, email FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DELETE
func (r *UserRepository) Delete(id int) error {
	result, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

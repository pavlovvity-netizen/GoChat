package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"GoChat/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(reg *domain.User) error {
	query := "INSERT INTO users (name, password_hash) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(query, reg.Username, reg.PasswordHash).Scan(&reg.ID)
	if err != nil {
		return fmt.Errorf("Ошибка при регистрации пользователя: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByName(name string) (*domain.User, error) {
	query := "SELECT id, name, password_hash FROM users WHERE name = $1"

	user := &domain.User{}
	err := r.db.QueryRow(query, name).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

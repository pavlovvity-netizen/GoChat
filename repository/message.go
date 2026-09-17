package repository

import (
	"database/sql"
	"fmt"

	"GoChat/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll() ([]domain.Message, error) {
	query := `
			SELECT m.id, m.text, m.post_time, m.user_id, COALESCE(u.name, '') AS username
			FROM messages m
			JOIN users u ON m.user_id = u.id
			ORDER BY m.post_time ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message

	for rows.Next() {
		var msg domain.Message

		err := rows.Scan(&msg.ID, &msg.Text, &msg.PostTime, &msg.UserID, &msg.Username)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *PostgresRepository) Create(msg *domain.Message) error {
	query := "INSERT INTO messages (text, post_time, user_id) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(query, msg.Text, msg.PostTime, msg.UserID).Scan(&msg.ID)
	if err != nil {
		return fmt.Errorf("ошибка создания сообщения: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Update(msg *domain.Message) error {
	query := `
		UPDATE messages 
		SET text = $1 
		WHERE id = $2 
		RETURNING post_time, user_id`

	res, err := r.db.Exec(query, msg.Text, msg.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления сообщения: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка обновления сообщения: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("сообщение не найдено: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *PostgresRepository) Delete(id int) error {
	query := "DELETE FROM messages WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

package repository

import (
	"family-tree-api/internal/database"
	"family-tree-api/internal/models"
	"fmt"
)

func CreateContactMessage(msg *models.ContactMessage) error {
	query := `
		INSERT INTO contact_messages (name, email, message, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`

	err := database.DB.QueryRow(query, msg.Name, msg.Email, msg.Message).
		Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func GetAllMessages() ([]models.ContactMessage, error) {
	query := `
		SELECT id, name, email, message, created_at
		FROM contact_messages ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.ContactMessage
	for rows.Next() {
		var msg models.ContactMessage
		err := rows.Scan(&msg.ID, &msg.Name, &msg.Email, &msg.Message, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

package repository

import (
	"database/sql"
	"errors"
	"family-tree-api/internal/database"
	"family-tree-api/internal/models"
	"fmt"
)

var ErrPasswordResetTokenNotFound = errors.New("password reset token not found")

func CreatePasswordResetToken(token *models.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`

	err := database.DB.QueryRow(query, token.UserID, token.Token, token.ExpiresAt).
		Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}
	return nil
}

func GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	resetToken := &models.PasswordResetToken{}
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM password_reset_tokens WHERE token = $1
	`

	err := database.DB.QueryRow(query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrPasswordResetTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}
	return resetToken, nil
}

func DeletePasswordResetToken(token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := database.DB.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}
	return nil
}

func DeleteExpiredPasswordResetTokens() error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := database.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired password reset tokens: %w", err)
	}
	return nil
}

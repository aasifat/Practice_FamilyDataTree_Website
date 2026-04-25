package repository

import (
	"database/sql"
	"errors"
	"family-tree-api/internal/database"
	"family-tree-api/internal/models"
	"fmt"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password, first_name, last_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	err := database.DB.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, "user").
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users WHERE email = $1
	`

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func UpdatePassword(userID int, hashedPassword string) error {
	query := `
		UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2
	`

	_, err := database.DB.Exec(query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := database.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// StoreOTP saves OTP to existing user
func StoreOTP(email, otp string, expiryMinutes int) error {
	query := `
		UPDATE users SET otp_code = $1, otp_expiry = CURRENT_TIMESTAMP + INTERVAL '1 minute' * $2
		WHERE email = $3
	`

	result, err := database.DB.Exec(query, otp, expiryMinutes, email)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with email %s not found", email)
	}

	return nil
}

// VerifyOTP checks if OTP is valid and not expired
func VerifyOTP(email, otp string) (bool, error) {
	var storedOTP sql.NullString
	var expiryTime sql.NullTime

	query := `
		SELECT otp_code, otp_expiry FROM users WHERE email = $1
	`

	err := database.DB.QueryRow(query, email).Scan(&storedOTP, &expiryTime)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}

	// Check if OTP exists and matches
	if !storedOTP.Valid || storedOTP.String != otp {
		return false, nil
	}

	// Check if OTP has expired
	if !expiryTime.Valid {
		return false, nil
	}

	// Compare with current time
	now := time.Now()
	if expiryTime.Time.Before(now) {
		return false, nil
	}

	return true, nil
}

// ClearOTP clears OTP after successful verification
func ClearOTP(email string) error {
	query := `
		UPDATE users SET otp_code = NULL, otp_expiry = NULL, is_verified = TRUE
		WHERE email = $1
	`

	_, err := database.DB.Exec(query, email)
	if err != nil {
		return fmt.Errorf("failed to clear OTP: %w", err)
	}
	return nil
}

// VerifyAndActivateUser marks an unverified user as verified and clears OTP data
func VerifyAndActivateUser(email string) error {
	query := `
		UPDATE users SET otp_code = NULL, otp_expiry = NULL, is_verified = TRUE
		WHERE email = $1 AND is_verified = FALSE
	`

	result, err := database.DB.Exec(query, email)
	if err != nil {
		return fmt.Errorf("failed to verify and activate user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already verified")
	}

	return nil
}

// StoreOTPWithSignupData stores OTP and signup data for verification
func StoreOTPWithSignupData(email, otp, hashedPassword, firstName, lastName string, expiryMinutes int) error {
	// First check if email already exists as a verified user
	var existingID int
	var isVerified bool
	checkQuery := `SELECT id, is_verified FROM users WHERE email = $1`
	err := database.DB.QueryRow(checkQuery, email).Scan(&existingID, &isVerified)

	if err == nil && isVerified {
		// User already exists and is verified
		return fmt.Errorf("email already registered")
	} else if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	// If unverified user exists, delete the old record first to avoid conflicts
	if err == nil && !isVerified {
		deleteQuery := `DELETE FROM users WHERE email = $1 AND is_verified = FALSE`
		_, err = database.DB.Exec(deleteQuery, email)
		if err != nil {
			return fmt.Errorf("failed to delete old signup record: %w", err)
		}
	}

	// Insert new signup record
	query := `
		INSERT INTO users (email, password, otp_code, otp_expiry, first_name, last_name, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP + INTERVAL '1 minute' * $4, $5, $6, 'user', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = database.DB.Exec(query, email, hashedPassword, otp, expiryMinutes, firstName, lastName)
	if err != nil {
		return fmt.Errorf("failed to store OTP with signup data: %w", err)
	}
	return nil
}

// GetSignupData retrieves signup data for a user
func GetSignupData(email string) (*models.SignupData, error) {
	signupData := &models.SignupData{}
	query := `
		SELECT email, password, first_name, last_name
		FROM users WHERE email = $1 AND is_verified = FALSE
	`

	err := database.DB.QueryRow(query, email).Scan(
		&signupData.Email, &signupData.HashedPassword, &signupData.FirstName, &signupData.LastName,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get signup data: %w", err)
	}

	// Check if password exists
	if signupData.HashedPassword == "" {
		return nil, fmt.Errorf("signup data is not available")
	}

	return signupData, nil
}

// ClearOTPAndSignupData clears OTP and signup data after successful signup
func ClearOTPAndSignupData(email string) error {
	query := `
		UPDATE users SET 
			otp_code = NULL, 
			otp_expiry = NULL, 
			signup_password = NULL,
			is_verified = TRUE
		WHERE email = $1
	`

	_, err := database.DB.Exec(query, email)
	if err != nil {
		return fmt.Errorf("failed to clear OTP and signup data: %w", err)
	}
	return nil
}

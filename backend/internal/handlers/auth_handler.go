package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"family-tree-api/config"
	"family-tree-api/internal/auth"
	"family-tree-api/internal/models"
	"family-tree-api/internal/repository"
	"family-tree-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// SignUp is deprecated - use RequestOTPForSignup + VerifyOTPAndSignup instead
// Kept for backward compatibility
func SignUp(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "direct signup is disabled. Use OTP-based signup: RequestOTPForSignup -> VerifyOTPAndSignup",
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !auth.VerifyPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		// Do not reveal whether the email exists.
		c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset link sent"})
		return
	}

	resetToken, err := generateSecureToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate reset token"})
		return
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	tokenRecord := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: expiresAt,
	}

	if err := repository.CreatePasswordResetToken(tokenRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reset token"})
		return
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.AppConfig.FrontendURL, resetToken)
	if err := utils.SendPasswordResetEmail(req.Email, resetURL); err != nil {
		// Log email sending failures internally, but always return a generic message.
		fmt.Printf("failed to send password reset email: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset link sent"})
}

func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenRecord, err := repository.GetPasswordResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	if tokenRecord.ExpiresAt.Before(time.Now()) {
		repository.DeletePasswordResetToken(req.Token)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	if err := repository.UpdatePassword(tokenRecord.UserID, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	if err := repository.DeletePasswordResetToken(req.Token); err != nil {
		fmt.Printf("failed to delete used reset token: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func generateSecureToken(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := repository.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	// Remove passwords
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

func DeleteUserByAdmin(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := repository.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// RequestOTPForSignup sends OTP for signup verification
func RequestOTPForSignup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate OTP
	otp := utils.GenerateOTP()

	// Hash password before storing
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	// Store OTP with signup data in database (expires in 10 minutes)
	if err := repository.StoreOTPWithSignupData(req.Email, otp, hashedPassword, req.FirstName, req.LastName, 10); err != nil {
		fmt.Printf("failed to store OTP with signup data: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send OTP via email
	emailErr := utils.SendOTPEmail(req.Email, otp)
	if emailErr != nil {
		fmt.Printf("[ERROR] Failed to send OTP email to %s: %v\n", req.Email, emailErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP email. " + emailErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email", "email": req.Email})
}

// VerifyOTPAndSignup verifies OTP and creates user account
func VerifyOTPAndSignup(c *gin.Context) {
	var req models.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if req.Email == "" || req.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and OTP are required"})
		return
	}

	// Verify OTP
	isValid, err := repository.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		fmt.Printf("[ERROR] OTP verification failed for %s: %v\n", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or OTP"})
		return
	}

	if !isValid {
		fmt.Printf("[INFO] Invalid or expired OTP for %s\n", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired OTP"})
		return
	}

	// Verify the user and mark as verified (user already exists, just needs to be activated)
	if err := repository.VerifyAndActivateUser(req.Email); err != nil {
		fmt.Printf("[ERROR] Failed to activate user %s: %v\n", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to activate user"})
		return
	}

	// Get the user with updated data
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		fmt.Printf("[ERROR] Failed to retrieve user %s: %v\n", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "signup successful",
		"token":   token,
		"user":    user,
	})
}

// RequestOTP sends OTP to the user's email
func RequestOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		// Do not reveal whether the email exists
		c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email if it exists"})
		return
	}

	// Generate OTP
	otp := utils.GenerateOTP()

	// Store OTP in database (expires in 10 minutes)
	if err := repository.StoreOTP(req.Email, otp, 10); err != nil {
		// User not found - should not happen, but handle gracefully
		c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email if it exists"})
		return
	}

	// Send OTP via email
	emailErr := utils.SendOTPEmail(req.Email, otp)
	if emailErr != nil {
		fmt.Printf("[ERROR] Failed to send OTP email to %s: %v\n", req.Email, emailErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP email. " + emailErr.Error()})
		return
	}

	_ = user // Avoid unused variable warning
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email", "email": req.Email})
}

// VerifyOTPAndLogin verifies OTP and logs in user
func VerifyOTPAndLogin(c *gin.Context) {
	var req models.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify OTP
	isValid, err := repository.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or OTP"})
		return
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired OTP"})
		return
	}

	// Get user
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Clear OTP after successful verification
	if err := repository.ClearOTP(req.Email); err != nil {
		fmt.Printf("failed to clear OTP: %v\n", err)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user":    user,
	})
}

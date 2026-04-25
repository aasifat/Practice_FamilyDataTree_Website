package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"

	"family-tree-api/config"
)

func SendPasswordResetEmail(toEmail, resetURL string) error {
	if config.AppConfig.SMTPHost == "" || config.AppConfig.SMTPUser == "" || config.AppConfig.SMTPPass == "" {
		return fmt.Errorf("smtp configuration is not set")
	}

	from := config.AppConfig.SMTPUser
	subject := "Password Reset Request"
	body := fmt.Sprintf("Hello,\n\nWe received a request to reset your password. Use the link below to reset it:\n\n%s\n\nIf you did not request this, you can safely ignore this email.\n\nThanks,\nFamily Tree Team", resetURL)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.SMTPHost, config.AppConfig.SMTPPort)
	auth := smtp.PlainAuth("", config.AppConfig.SMTPUser, config.AppConfig.SMTPPass, config.AppConfig.SMTPHost)

	return smtp.SendMail(serverAddr, auth, from, []string{toEmail}, []byte(msg))
}

// SendOTPEmail sends an OTP to the user's email
func SendOTPEmail(toEmail, otp string) error {
	if config.AppConfig.SMTPHost == "" || config.AppConfig.SMTPUser == "" || config.AppConfig.SMTPPass == "" {
		return fmt.Errorf("SMTP is not configured. Please set SMTP_HOST, SMTP_USER, and SMTP_PASS environment variables. To test without email, check server logs for OTP")
	}

	from := config.AppConfig.SMTPUser
	subject := "Your OTP for Family Tree Login"
	body := fmt.Sprintf("Hello,\n\nYour One-Time Password (OTP) for logging into Family Tree is:\n\n%s\n\nThis OTP will expire in 10 minutes.\n\nIf you did not request this, please ignore this email.\n\nThanks,\nFamily Tree Team", otp)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.SMTPHost, config.AppConfig.SMTPPort)
	auth := smtp.PlainAuth("", config.AppConfig.SMTPUser, config.AppConfig.SMTPPass, config.AppConfig.SMTPHost)

	if err := smtp.SendMail(serverAddr, auth, from, []string{toEmail}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send OTP email to %s: %w", toEmail, err)
	}
	
	return nil
}

// GenerateOTP generates a random 6-digit OTP
func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

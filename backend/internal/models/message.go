package models

import "time"

type ContactMessage struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Message   string    `json:"message" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required,min=10"`
}

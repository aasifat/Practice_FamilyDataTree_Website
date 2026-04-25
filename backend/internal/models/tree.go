package models

import "time"

type FamilyTree struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Members   []Person  `json:"members,omitempty"`
}

type TreeCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type TreeUpdateRequest struct {
	Name string `json:"name"`
}

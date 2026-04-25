package models

import "time"

type Person struct {
	ID        int       `json:"id"`
	TreeID    int       `json:"tree_id"`
	Name      string    `json:"name" binding:"required"`
	Gender    string    `json:"gender"`    // "male" or "female"
	FatherID  *int      `json:"father_id"` // NULL for root/maternal lineage
	MotherID  *int      `json:"mother_id"` // NULL for root/paternal lineage
	SpouseID  *int      `json:"spouse_id"` // Optional spouse link (bidirectional)
	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Children  []Person  `json:"children,omitempty"`
	Spouse    *Person   `json:"spouse,omitempty"`
}

type PersonRequest struct {
	Name     string `json:"name" binding:"required"`
	Gender   string `json:"gender" binding:"required"` // "male" or "female"
	FatherID *int   `json:"father_id"`
	MotherID *int   `json:"mother_id"`
	SpouseID *int   `json:"spouse_id"`
}

type PersonUpdateRequest struct {
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	SpouseID *int   `json:"spouse_id"`
}

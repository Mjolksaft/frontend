package structs

import (
	"time"
)

type Password struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"createAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	HashedPassword string    `json:"hashedPassword"`
	Application    string    `json:"application"`
}

type User struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	HashedPassword string    `json:"hashed_password"`
	Username       string    `json:"username"`
	IsAdmin        bool      `json:"is_admin"`
}

package structs

import (
	"time"
)

type Password struct {
	ID             string `json:"ID"`
	HashedPassword string `json:"HashedPassword"`
	Application    string `json:"ApplicationName"`
	UserID         string `json:"UserID"`
}

type User struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	HashedPassword string    `json:"hashed_password"`
	Username       string    `json:"username"`
	IsAdmin        bool      `json:"is_admin"`
}

type Command struct {
	Name        string
	Description string
	Callback    func([]string, MenuSwitcher) error
}

type Menu struct {
	Prefix   string
	Commands map[string]Command
}

// MenuSwitcher interface defines the methods for switching and getting menus
type MenuSwitcher interface {
	SwitchMenu(int)
	GetCurrentMenu() Menu
}

type GetUserInfo struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Username     string    `json:"username"`
	Applications []string  `json:"applications"`
}

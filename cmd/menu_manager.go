package cmd

import "frontend/structs"

// MenuSwitcher interface defines the methods for switching and getting menus
type MenuSwitcher interface {
	SwitchMenu(int)
	GetCurrentMenu() structs.Menu
}

// MenuManager struct implements the MenuSwitcher interface
type MenuManager struct {
	Menus       []structs.Menu
	CurrentMenu int
}

func NewMenuManager() *MenuManager {
	return &MenuManager{
		Menus:       []structs.Menu{}, // Menus will be initialized later
		CurrentMenu: MainMenu,         // Start with the main menu
	}
}

// 0 mainMenu, 1 vault
func (m *MenuManager) SwitchMenu(menu int) {
	m.CurrentMenu = menu
}

func (m *MenuManager) GetCurrentMenu() structs.Menu {
	return m.Menus[m.CurrentMenu]
}

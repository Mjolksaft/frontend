package cmd

import (
	"bufio"
	"fmt"
	"frontend/internal/commands"
	"frontend/structs"
	"os"
	"strings"
)

const (
	MainMenu int = iota
	VaultMenu
)

func CLILoop(m *MenuManager) {
	reader := bufio.NewReader(os.Stdin)

	// runs the login command on fail close program
	command := m.GetCurrentMenu().Commands["login"]
	err := command.Callback(nil, m)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		menu := m.GetCurrentMenu()

		// write the prefix
		fmt.Print(menu.Prefix)
		// grab the input
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading string: %w", err)
			continue
		}

		// clean it
		cleanData := strings.TrimSpace(data)
		splitData := strings.Split(cleanData, " ")

		// grab the command
		command, exists := menu.Commands[cleanData]
		if !exists {
			fmt.Println("option does not exist")
			continue
		}

		// run the function
		err = command.Callback(splitData, m)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func MenuInit(m *MenuManager) { // make the help and back commands global commands (not a specific menu option) that takes a menu to back to and help takes the current menu
	m.Menus = []structs.Menu{
		{Prefix: "main > ",
			Commands: map[string]structs.Command{
				"help": {
					Name:        "Help",
					Description: "Shows menu options",
					Callback:    commands.HelpCommand,
				},
				"exit": {
					Name:        "Exit",
					Description: "Exits the program",
					Callback:    commands.ExitCommand,
				},
				"login": {
					Name:        "login",
					Description: "logs in new user",
					Callback:    commands.LoginCommand,
				},
				"vault": {
					Name:        "Vault",
					Description: "enter the vault",
					Callback:    commands.EnterVault,
				},
			},
		},
		{Prefix: "vault > ",
			Commands: map[string]structs.Command{
				"help": {
					Name:        "Help",
					Description: "Shows menu options",
					Callback:    commands.HelpCommand,
				},
				"back": {
					Name:        "Back",
					Description: "Backs to main menu",
					Callback:    commands.BackCommand,
				},
				"create": {
					Name:        "Create",
					Description: "Create a password",
					Callback:    commands.CreatePasswordCommand,
				},
				"get": {
					Name:        "Get",
					Description: "Get a password by application",
					Callback:    commands.GetPasswordByApplicationCommand,
				},
				"get_all": {
					Name:        "Get",
					Description: "Get a password by application",
					Callback:    commands.GetPasswordsCommand,
				},
				"delete": {
					Name:        "Delete",
					Description: "Delete by application",
					Callback:    commands.DeletePasswordCommand,
				},
				"profile": {
					Name:        "Profile",
					Description: "Get user info",
					Callback:    commands.GetUserInfo,
				},
			},
		},
	}
}

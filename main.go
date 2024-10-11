package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"frontend/structs"
	"net/http"
	"os"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Callback    func([]string) error
}

type Menu struct {
	Prefix   string
	Commands map[string]Command
}

var Menus []Menu

func menuInit() { // make the help and back commands global commands (not a specific menu option) that takes a menu to back to and help takes the current menu
	Menus = []Menu{
		{Prefix: "main > ",
			Commands: map[string]Command{
				"help": {
					Name:        "Help",
					Description: "Shows menu options",
					Callback:    helpCommand,
				},
				"exit": {
					Name:        "Exit",
					Description: "Exits the program",
					Callback:    exitCommand,
				},
				"login": {
					Name:        "Login",
					Description: "login user",
					Callback:    loginCommand,
				},
				"vault": {
					Name:        "Vault",
					Description: "enter the vault",
					Callback:    enterVault,
				},
			},
		},
		{Prefix: "vault > ",
			Commands: map[string]Command{
				"help": {
					Name:        "Help",
					Description: "Shows menu options",
					Callback:    helpCommand,
				},
				"back": {
					Name:        "Back",
					Description: "Backs to main menu",
					Callback:    backCommand,
				},
				"create": {
					Name:        "Create",
					Description: "Create a password",
					Callback:    createPassword,
				},
				"update": {
					Name:        "Update",
					Description: "Update a password",
					Callback:    updatePasswordCommand,
				},
				"delete": {
					Name:        "Delete",
					Description: "Delete a password",
					Callback:    deletePasswordCommand,
				},
			},
		},
	}
}

const (
	mainMenu int = iota
	vaultMenu
)

var currentMenu = mainMenu

func main() {
	menuInit()
	CLILoop(mainMenu)
}

func CLILoop(menuOption int) {
	// set the menu
	currentMenu = menuOption
	menu := Menus[currentMenu]

	reader := bufio.NewReader(os.Stdin)
	for {
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
		err = command.Callback(splitData)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func helpCommand(arguments []string) error {
	i := 1
	for key, value := range Menus[currentMenu].Commands {
		fmt.Printf("%d: %s (%s)\n", i, key, value.Description)
		i++
	}

	return nil
}

func exitCommand(arguments []string) error {
	fmt.Println("exiting program")
	os.Exit(0)

	return nil
}

func loginCommand(arguments []string) error {
	fmt.Println("Login")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("username >")
	username, err := reader.ReadString('\n')

	if err != nil {
		return fmt.Errorf("error reading username: %w", err)
	}
	fmt.Print("password >")
	password, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading username: %w", err)
	}
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	jsonString := fmt.Sprintf(`{"password": "%s", "username": "%s"}`, password, username)

	ioReader := strings.NewReader(jsonString)

	// Make the request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/login", ioReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error doing request %w", err)
	}

	if res.StatusCode == 200 { // if the request is not good send error
		decoder := json.NewDecoder(res.Body)

		var data structs.User
		if err := decoder.Decode(&data); err != nil {
			return fmt.Errorf("error decoding data: %w", err)
		}

		fmt.Println(data)
	}

	return nil
}

func createPassword(arguments []string) error {
	// take the input of password then the application
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("password > ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("erro reading input: %w", err)
	}

	trimedPassword := strings.TrimSpace(input)

	fmt.Print("application > ")
	input, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("erro reading input: %w", err)
	}

	trimedApplication := strings.TrimSpace(input)

	// create json string
	jsonString := fmt.Sprintf(`{
			"password": "%s",
			"application": "%s"
		}`,
		trimedPassword,
		trimedApplication,
	)

	// Send the POST request
	ioReader := strings.NewReader(jsonString)
	res, err := http.Post("http://localhost:8080/api/passwords", "application/json", ioReader)
	if err != nil {
		return fmt.Errorf("post request error: %w", err)
	}
	defer res.Body.Close()

	// read the response
	fmt.Println(res)

	return nil
}

func updatePasswordCommand(args []string) error {
	fmt.Println("update password")
	return nil
}

func deletePasswordCommand(args []string) error {
	fmt.Println("delete password")
	return nil
}

func enterVault(args []string) error {
	fmt.Println("enter vault menu")
	CLILoop(vaultMenu)
	return nil
}

func backCommand(args []string) error {
	fmt.Println("back to main menu")
	CLILoop(mainMenu)
	return nil
}

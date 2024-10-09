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

var Menus []map[string]Command

func menuInit() {
	Menus = []map[string]Command{
		{
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
		},
	}

}

const (
	mainMenu int = iota
	secondMenu
)

var currentMenu = mainMenu

func main() {
	menuInit()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		data, err := reader.ReadString('\n')
		cleanData := strings.TrimSpace(data)
		if err != nil {
			fmt.Println("error reading string: %w", err)
			continue
		}
		menu := Menus[currentMenu]
		command, exists := menu[cleanData]
		if !exists {
			fmt.Println("option does not exist")
		}

		splitData := strings.Split(cleanData, " ")
		err = command.Callback(splitData)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func helpCommand(arguments []string) error {
	i := 1
	for key, value := range Menus[currentMenu] {
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

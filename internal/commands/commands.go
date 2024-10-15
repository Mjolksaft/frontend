package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"frontend/structs"
	"net/http"
	"os"
	"strings"
)

func HelpCommand(arguments []string, m structs.MenuSwitcher) error {
	i := 1
	for key, value := range m.GetCurrentMenu().Commands {
		fmt.Printf("%d: %s (%s)\n", i, key, value.Description)
		i++
	}

	return nil
}

func ExitCommand(arguments []string, m structs.MenuSwitcher) error {
	fmt.Println("exiting program")
	os.Exit(0)

	return nil
}

func LoginCommand(arguments []string, m structs.MenuSwitcher) error {
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

func CreatePassword(arguments []string, m structs.MenuSwitcher) error {
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

func UpdatePasswordCommand(args []string, m structs.MenuSwitcher) error {
	fmt.Println("update password")
	return nil
}

func DeletePasswordCommand(args []string, m structs.MenuSwitcher) error {
	fmt.Println("delete password")
	return nil
}

func EnterVault(args []string, m structs.MenuSwitcher) error {
	fmt.Println("enter vault menu")
	m.SwitchMenu(1)
	return nil
}

func BackCommand(args []string, m structs.MenuSwitcher) error {
	fmt.Println("back to main menu")
	m.SwitchMenu(0)
	return nil
}

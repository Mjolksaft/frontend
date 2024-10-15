package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"frontend/structs"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var client *http.Client

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

	// Get username input
	fmt.Print("username >")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading username: %w", err)
	}

	// Get password input
	fmt.Print("password >")
	password, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading password: %w", err)
	}

	// Trim spaces from inputs
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Create JSON payload for login
	jsonString := fmt.Sprintf(`{"password": "%s", "username": "%s"}`, password, username)
	ioReader := strings.NewReader(jsonString)

	// Define the login URL
	loginURL := "http://localhost:8080/api/login"

	// Create the HTTP request
	req, err := http.NewRequest("POST", loginURL, ioReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("error creating cookie jar: %w", err)
	}

	// Create an HTTP client with the cookie jar
	client = &http.Client{
		Jar: jar,
	}

	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close() // Ensure the response body is closed

	// Parse the URL to retrieve cookies
	u, err := url.Parse(loginURL)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	// Print cookies stored in the cookie jar after login
	cookies := jar.Cookies(u)
	fmt.Println("Cookies after login:", cookies)

	// Check if the login was successful
	if res.StatusCode == 200 {
		decoder := json.NewDecoder(res.Body)

		// Decode the response body into the User struct
		var data structs.User
		if err := decoder.Decode(&data); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

	} else {
		return fmt.Errorf("login failed with status code: %d", res.StatusCode)
	}

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

func TestCommand(args []string, m structs.MenuSwitcher) error {
	fmt.Println("Test the api")

	_, err := client.Get("http://localhost:8080/api/users")
	if err != nil {
		return fmt.Errorf("error with request: %w", err)
	}

	// decode res
	return nil
}

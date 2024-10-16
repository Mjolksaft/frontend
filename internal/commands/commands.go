package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"frontend/internal/encryption"
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

	input, err := getInput([]string{"username", "password"})
	if err != nil {
		return fmt.Errorf("error recieving input: %w", err)
	}

	// Create JSON payload for login
	jsonString := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, input[0], input[1])
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

func TestEncryption(args []string, m structs.MenuSwitcher) error {
	input, err := getInput([]string{"master password", "password"})
	if err != nil {
		return fmt.Errorf("error recieving input: %w", err)
	}

	encryptedPassword, err := encryption.EncryptPassword(input[0], input[1])
	if err != nil {
		return fmt.Errorf("error encrypting password: %w", err)
	}

	password, err := encryption.DecryptPassword(input[0], encryptedPassword)
	if err != nil {
		return fmt.Errorf("error decrypting password: %w", err)
	}

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("encrypted: %s\n", encryptedPassword)

	return nil

}

func getInput(queries []string) ([]string, error) {
	length := len(queries)
	var input = make([]string, length)

	reader := bufio.NewReader(os.Stdin)
	for i, query := range queries {
		fmt.Printf("%s > ", query)
		value, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading password: %w", err)
		}

		input[i] = strings.TrimSpace(value)
	}
	return input, nil
}

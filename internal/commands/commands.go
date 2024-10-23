package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"frontend/internal/encryption"
	"frontend/internal/util"
	"frontend/structs"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/atotto/clipboard"
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

func CreatePasswordCommand(args []string, m structs.MenuSwitcher) error {
	// check if client exists to see if the user is logged in or not

	// get the input
	input, err := getInput([]string{"masterPassword", "password"})
	if err != nil {
		return fmt.Errorf("error getting input: %w", err)
	}

	// wait of screen change to get the application you want
	windowTitle := util.MonitorWindowChange()
	encodedAppName := url.QueryEscape(strings.Split(windowTitle, " - ")[0])

	// encrypt the password
	encrypted, err := encryption.EncryptPassword(input[0], input[1])
	if err != nil {
		return fmt.Errorf("error encrypting: %w", err)
	}

	//make the json string
	jsonString := fmt.Sprintf(`{"password": "%s", "application": "%s"}`, encrypted, encodedAppName)
	reader := strings.NewReader(jsonString)

	// create the request
	res, err := client.Post("http://localhost:8080/api/passwords", "text/plain", reader)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// check the results
	if res.StatusCode == 201 {
		fmt.Printf("%s password added\n", encodedAppName)
		return nil
	}

	return nil
}

func GetPasswordCommand(args []string, m structs.MenuSwitcher) error {
	type Password struct {
		HashedPassword  string `json:"HashedPassword"`
		ApplicationName string `json:"ApplicationName"`
	}

	// query the masterpassword for decryption
	input, err := getInput([]string{"master password"})
	if err != nil {
		return fmt.Errorf("error getting input: %w", err)
	}

	// wait of screen change to get the application you want
	windowTitle := util.MonitorWindowChange()
	encodedAppName := url.QueryEscape(strings.Split(windowTitle, " - ")[0])

	// make the request
	fullUrl := fmt.Sprintf("http://localhost:8080/api/passwords?application_name=%s", encodedAppName)
	fmt.Println(fullUrl)
	res, err := client.Get(fullUrl)
	if err != nil {
		return fmt.Errorf("error with request: %w", err)
	}

	// decode the password
	decoder := json.NewDecoder(res.Body)
	var body Password
	if err := decoder.Decode(&body); err != nil {
		return fmt.Errorf("error decoding body: %w", err)
	}

	// decrypt the password
	password, err := encryption.DecryptPassword(input[0], body.HashedPassword)
	if err != nil {
		return fmt.Errorf("error decrypting password: %w", err)
	}

	// add to clip board
	err = clipboard.WriteAll(password)
	if err != nil {
		return err
	}

	fmt.Println("Text copied to clipboard!")

	return nil
}

func GetPasswordsCommand(args []string, m structs.MenuSwitcher) error {
	type passwordStruct struct {
		ApplicationName string `json:"applicationName"`
		HashedPassword  string `json:"hashedPassword"`
	}

	// make the request
	res, err := client.Get("http://localhost:8080/api/passwords")
	if err != nil {
		return fmt.Errorf("error with request: %w", err)
	}

	// decode the password
	decoder := json.NewDecoder(res.Body)
	var body []passwordStruct
	if err := decoder.Decode(&body); err != nil {
		return fmt.Errorf("error decoding body: %w", err)
	}

	for i, data := range body {
		fmt.Printf("%d. %s \n", i+1, data.ApplicationName)

	}

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

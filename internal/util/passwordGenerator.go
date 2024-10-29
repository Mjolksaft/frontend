package util

import (
	"fmt"
	"math/rand"
)

// creates a random password of set length which can include synbols or numbers
func GeneratePassword(length int, symbols, numbers bool) string {
	password := make([]byte, length)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	if symbols {
		charset += "!@#$%^&*()_+=-"
	}
	if numbers {
		charset += "0123456789"
	}
	fmt.Println(charset)
	for i := range password {
		randNum := rand.Intn(len(charset))
		password[i] = charset[randNum]
	}

	return string(password)
}

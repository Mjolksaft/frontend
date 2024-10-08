package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		data, err := reader.ReadString('\n')
		cleanData := strings.TrimSpace(data)
		if err != nil {
			fmt.Println("error reading string: %w", err)
			continue
		}
		fmt.Println("hello world")
	}

}

package main

import (
	"frontend/cmd"
)

func main() {
	menuManager := cmd.NewMenuManager()
	cmd.MenuInit(menuManager)
	cmd.CLILoop(menuManager)
}

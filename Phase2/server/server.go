package main

import (
	"phase2/server/commands"
	"phase2/server/console"
)

func main() {
	console.Con.App.ActiveMenu().SetCommands(commands.Commands())
	console.Con.App.Menu("session").SetCommands(commands.SessionCmd())
	console.Con.App.Start()
}

package main

import (
	"phase3/server/commands"
	"phase3/server/console"
)

func main() {
	console.Con.App.ActiveMenu().SetCommands(commands.Commands())
	console.Con.App.Menu("session").SetCommands(commands.SessionCmd())
	console.Con.App.Start()
}

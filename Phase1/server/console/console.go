package console

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/reeflective/console"
	"github.com/reeflective/readline"
)

var Con *ServerConsole

type ServerConsole struct {
	App        *console.Console
	Jobs       map[int]Job
	JobID      int
	Sessions   map[int]*Session
	SessionID  int
	SessionAct int
	PrintWait  chan bool
}

type Job struct {
	ID      int
	Name    string
	Target  string
	JobCtrl chan bool
}

type Session struct {
	ID          int
	Target      string
	System      string
	Username    string
	Hostname    string
	SessionCtrl chan bool
	Conn        net.Conn
}

func init() {
	Con = &ServerConsole{
		App:        console.New("main"),
		Jobs:       make(map[int]Job),
		JobID:      0,
		Sessions:   make(map[int]*Session),
		SessionID:  0,
		SessionAct: 0,
		PrintWait:  make(chan bool),
	}

	app := Con.App
	app.NewlineBefore = true
	app.NewlineAfter = true
	app.SetPrintLogo(func(_ *console.Console) {
		fmt.Print(`
A simple logo............ 	
`)
	})

	menu := app.ActiveMenu()
	setupPrompt(menu)
	menu.AddInterrupt(readline.ErrInterrupt, exitConfirm)

	// session menu
	s := app.NewMenu("session")
	sessionPrompt(s)
}

func (con *ServerConsole) NextJobID() int {
	next := con.JobID + 1
	con.JobID++
	return next
}

func (con *ServerConsole) NextSessionID() int {
	next := con.SessionID + 1
	con.SessionID++
	return next
}

// setupPrompt is a function which sets up the prompts for the main menu.
func setupPrompt(m *console.Menu) {
	p := m.Prompt()

	p.Primary = func() string {
		prompt := "[main] > "
		return prompt
	}
}

func sessionPrompt(m *console.Menu) {

	p := m.Prompt()
	p.Primary = func() string {
		// Bold Red Normal
		return fmt.Sprintf("%s%ssession %s> ", "\033[1m", "\033[31m", "\033[0m")
	}

	p.Secondary = func() string {
		return "> "
	}
}

// Used for Ctrl-C exit
func exitConfirm(_ *console.Console) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (Y/y): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		os.Exit(0)
	}
}

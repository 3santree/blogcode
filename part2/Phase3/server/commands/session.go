package commands

import (
	"net"
	"phase3/server/console"

	"github.com/jedib0t/go-pretty/table"
)

func session(k, i int) {
	if i != 0 {
		PrintSuccess("Interact with Session %d\n", i)
		console.Con.SessionAct = i
		console.Con.App.SwitchMenu("session")
		return
	}

	if k == 0 {
		printSessions()
	} else {
		if session, ok := console.Con.Sessions[k]; ok {
			session.SessionCtrl <- true
		} else {
			PrintError("session %d not exist\n", k)
		}
	}
}

func printSessions() {
	if len(console.Con.Sessions) == 0 {
		PrintInfo("No Sessions\n")
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(BorderTable)
	tw.AppendHeader(table.Row{
		"ID",
		"Target",
		"OS",
		"Username",
		"Hostname",
	})

	for id := range console.Con.Sessions {
		tw.AppendRow(table.Row{
			console.Con.Sessions[id].ID,
			console.Con.Sessions[id].Target,
			console.Con.Sessions[id].System,
			console.Con.Sessions[id].Username,
			console.Con.Sessions[id].Hostname,
		})
	}
	PrintRaw("%s\n", tw.Render())
}

func newSession(conn net.Conn) *console.Session {
	id := console.Con.NextSessionID()
	console.Con.Sessions[id] = &console.Session{
		ID:          id,
		Target:      conn.RemoteAddr().String(),
		SessionCtrl: make(chan bool),
		Conn:        conn,
	}
	PrintEventSuccess("Session #%d from %s\n", id, conn.RemoteAddr().String())

	go func() {
		<-console.Con.Sessions[id].SessionCtrl
		conn.Close()
		delete(console.Con.Sessions, id)
	}()

	return console.Con.Sessions[id]
}

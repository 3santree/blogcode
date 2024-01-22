package commands

import (
	"fmt"
	"net"
	"phase2/server/console"

	"github.com/jedib0t/go-pretty/table"
)

func job(k int) {
	if k == 0 {
		jobPrint()
	} else {
		jobKill(k)
	}
}

func jobPrint() {
	if len(console.Con.Jobs) == 0 {
		fmt.Printf("No Jobs\n")
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(TableDefault)
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Target",
	})

	for id := range console.Con.Jobs {
		tw.AppendRow(table.Row{
			fmt.Sprintf("%d", id),
			console.Con.Jobs[id].Name,
			console.Con.Jobs[id].Target,
		})
	}
	fmt.Printf("%s\n", tw.Render())
}

func jobKill(id int) {
	console.Con.Jobs[id].JobCtrl <- true
	PrintInfo("Job #%d stopped (%s %s)\n", id, console.Con.Jobs[id].Name, console.Con.Jobs[id].Target)
}

func addJob(ip net.IP, port uint32, name string) {

	id := console.Con.NextJobID()
	console.Con.Jobs[id] = console.Job{
		ID:      id,
		Name:    name,
		Target:  fmt.Sprintf("%s:%d", ip.String(), port),
		JobCtrl: make(chan bool),
	}

	PrintSuccess("Successfully started job #%d\n", id)

	go func() {
		<-console.Con.Jobs[id].JobCtrl
		delete(console.Con.Jobs, id)
	}()
}

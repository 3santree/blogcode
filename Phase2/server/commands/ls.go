package commands

import (
	"phase2/pb"
	"phase2/server/console"
)

func ls() {

	conn := console.Con.Sessions[console.Con.SessionAct].Conn
	req := &pb.Envelope{
		ID:   1,
		Type: 1,
		Data: []byte("ls"),
	}
	err := WriteEnvelope(conn, req)
	if err != nil {
		PrintError("ls Error\n")
	}

	<-console.Con.PrintWait
}

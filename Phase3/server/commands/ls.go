package commands

import (
	"phase3/pb"
	"phase3/server/console"
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

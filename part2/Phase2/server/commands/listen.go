package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"phase2/pb"
	"phase2/server/console"

	"google.golang.org/protobuf/proto"
)

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024

	// ServerMaxMessageSize - Server-side max GRPC message size
	ServerMaxMessageSize = 2 * gb
)

func listen(ip net.IP, port uint32) {
	tcpListen(ip, port)
	addJob(ip, port, "tcp")
}

func tcpListen(ip net.IP, port uint32) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		fmt.Printf("listen err: %v\n", err)
		return
	}

	go acceptConnections(ln)
}

func acceptConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errType, ok := err.(*net.OpError); ok && errType.Op == "accept" {
				break // Listener was closed by the user
			}
			continue
		}
		if heartBeat(conn) {
			go handleConnection(conn)
		} else {
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		PrintEventError("Session closed for %s\n", conn.RemoteAddr())
	}()

	session := newSession(conn)
	message := make(chan *pb.Envelope)

	done := make(chan struct{})
	// Read loop
	go func() {
		for {
			envelope, err := ReadEnvelope(conn)
			if err != nil && !heartBeat(conn) {
				// fmt.Printf("Connection detected closed %v\n", err)
				session.SessionCtrl <- true
				done <- struct{}{}
				return
			}
			// echo message back
			// Here should be the message handler
			message <- envelope
		}
	}()
	// Send loop
Loop:
	for {
		select {
		case m := <-message:
			err := msgHandler(m, conn, session)
			if err != nil {
				break Loop
			}

		case <-done:
			break Loop

		}
	}
}

func WriteEnvelope(connection net.Conn, envelope *pb.Envelope) error {
	data, err := proto.Marshal(envelope)
	if err != nil {
		// fmt.Printf("Envelope marshaling error: %v\n", err)
		return err
	}
	dataLengthBuf := new(bytes.Buffer)
	binary.Write(dataLengthBuf, binary.LittleEndian, uint32(len(data)))
	connection.Write(dataLengthBuf.Bytes())
	connection.Write(data)
	return nil
}

func ReadEnvelope(connection net.Conn) (*pb.Envelope, error) {
	// Read the first four bytes to determine data length
	dataLengthBuf := make([]byte, 4) // Size of uint32
	n, err := io.ReadFull(connection, dataLengthBuf)
	if err != nil || n != 4 {
		// fmt.Printf("Socket error (read msg-length): %v\n", err)
		return nil, err
	}

	dataLength := int(binary.LittleEndian.Uint32(dataLengthBuf))
	if dataLength <= 0 || ServerMaxMessageSize < dataLength {
		// {{if .Config.Debug}}
		fmt.Printf("[pivot] read error: %s\n", err)
		// {{end}}
		return nil, errors.New("[pivot] invalid data length")
	}

	dataBuf := make([]byte, dataLength)

	n, err = io.ReadFull(connection, dataBuf)
	if err != nil || n != dataLength {
		fmt.Printf("Socket error (read data): %v\n", err)
		return nil, err
	}

	// Unmarshal the protobuf envelope
	envelope := &pb.Envelope{}
	err = proto.Unmarshal(dataBuf, envelope)
	if err != nil {
		fmt.Printf("Un-marshaling envelope error: %v\n", err)
		return nil, err
	}
	return envelope, nil
}

func heartBeat(connection net.Conn) bool {
	ping := &pb.Envelope{
		ID:   0,
		Type: 0,
		Data: []byte("ping"),
	}
	err := WriteEnvelope(connection, ping)
	if err != nil {
		fmt.Printf("heartBeat write failed %v\n", err)
		return false
	}

	read := make(chan *pb.Envelope)
	go func() {
		envelope, err := ReadEnvelope(connection)
		if err != nil {
			// fmt.Printf("heartBeat read failed %v\n", err)
			read <- &pb.Envelope{
				ID:   1,
				Type: 0,
				Data: []byte(""),
			}
		}
		read <- envelope
	}()

	select {
	case <-time.After(1 * time.Second):
		return false
	case m := <-read:
		if m.Type == 0 && m.ID == 0 {
			return true
		}
	}

	return false
}

func msgHandler(env *pb.Envelope, conn net.Conn, session *console.Session) error {
	err := error(nil)
	switch env.Type {
	// heartbeat
	case 0:
		echo := &pb.Envelope{
			ID:   1,
			Type: 0,
			Data: []byte("Server echo"),
		}
		err = WriteEnvelope(conn, echo)

	// cmd's response
	case 1:
		data := env.GetData()
		PrintInfo("%s", data)
		console.Con.PrintWait <- true
	// Receive session info
	case 2:
		data := env.GetData()
		info := &pb.SessionInfo{}
		err = proto.Unmarshal(data, info)
		if err != nil {
			fmt.Printf("Un-marshaling envelope error: %v\n", err)
		}
		session.System = info.GetSystem()
		session.Hostname = info.GetHostname()
		session.Username = info.GetUsername()
	}
	return err
}

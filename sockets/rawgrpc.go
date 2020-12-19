package sockets

import (
	"context"
	"log"
	"time"

	pb "github.com/preludeorg/pneuma/sockets/protos"
	"google.golang.org/grpc"
)

type GRPC struct {}

func init() {
	CommunicationChannels["grpc"] = GRPC{}
}

func (contact GRPC) Communicate(address string, sleep int, beacon Beacon) {
	for {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("[-] %s is either unavailable or a firewall is blocking traffic.", address)
		} else {
			interact(conn, beacon)
		}
		jitterSleep(sleep, "GRPC")
	}
}

func interact(conn *grpc.ClientConn, beacon Beacon) {
	defer conn.Close()
	c := pb.NewBeaconClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Handle(ctx, &pb.BeaconIncoming{Name: beacon.Name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetName())
}

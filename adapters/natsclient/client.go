// Package natsclient provides functionalities for interacting with a NATS
// broker for fetching input data.
package natsclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATSConnection holds the infrastructure concerns regarding the interaction
// with the NATS broker.
type NATSConnection struct {
	NATSConn  *nats.Conn
	Consumer  jetstream.Consumer
	Ctx       context.Context
	CancelCtx context.CancelFunc
}

// DashboardClient represents the domain-specific API that the service uses
// when interacting with the data source, hiding any infrastucture concerns.
type DashboardClient interface {
	FetchMessage() string
}

// DashboardNATSClient interacts with a NATS broker for fetching input data.
type DashboardNATSClient struct {
	conn *NATSConnection
}

// NewDashboardNATSClient is a constructor function for creating a
// connected and operational DashboardNATSClient.
func NewDashboardNATSClient() (*DashboardNATSClient, error) {
	natsConn := connectToNATS()
	return &DashboardNATSClient{conn: &natsConn}, nil
}

// FetchMessage fetches a single message from the NATS broker.
func (d *DashboardNATSClient) FetchMessage() []byte {
	messageData := ""
	msgs, err := d.conn.Consumer.Fetch(1)
	if err != nil {
		messageData += fmt.Sprintf("failed to fetch messages: %s\n", err)
		log.Print(messageData)
	}

	for msg := range msgs.Messages() {
		msg.Ack()
		messageData = string(msg.Data())
	}
	if msgs.Error() != nil {
		messageData += fmt.Sprintf("Error during Fetch(): %s\n", msgs.Error())
	}
	return []byte(messageData)
}

func connectToNATS() NATSConnection {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

	natsURI, found := os.LookupEnv("NATS_SERVER_URI")
	if !found {
		log.Fatal("ERROR: NATS server URI not set.")
	}
	log.Printf("The NATS URI is %s.", natsURI)

	nc, err := nats.Connect(natsURI)
	if err != nil {
		log.Fatalf("failed to connect to nats: %s.", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("ERROR: failed to create jetstream: %s.", err)
	}

	s, err := js.Stream(ctx, "TMA")
	if err != nil {
		log.Fatalf("ERROR: failed to initialize stream: %s.", err)
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "TMA",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("ERROR: failed to initialize consumer: %s.", err)
	}

	return NATSConnection{nc, c, ctx, cancel}
}

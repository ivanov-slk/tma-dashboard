// Package natsclient provides functionalities for interacting with a NATS
// broker for fetching input data.
package natsclient

import (
	"context"
	"fmt"
	"log"
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
	natsConn, err := connectToNATS()
	if err != nil {
		return nil, err
	}
	return &DashboardNATSClient{conn: &natsConn}, nil
}

// FetchMessage fetches a single message from the NATS broker.
func (d *DashboardNATSClient) FetchMessage() string {
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
	return messageData
}

func connectToNATS() (NATSConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

	nc, err := nats.Connect("http://nats-server:4222")
	if err != nil {
		log.Printf("failed to connect to nats: %s", err)
		return NATSConnection{nil, nil, ctx, cancel}, err
	}

	js, _ := jetstream.New(nc)

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TMA",
		Subjects: []string{"generated-data"},
	})
	if err != nil {
		return NATSConnection{nil, nil, ctx, cancel}, err
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "TMA",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return NATSConnection{nil, nil, ctx, cancel}, err
	}

	return NATSConnection{nc, c, ctx, cancel}, nil
}

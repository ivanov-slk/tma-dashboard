package natsclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSConnection struct {
	NATSConn  *nats.Conn
	Consumer  jetstream.Consumer
	Ctx       context.Context
	CancelCtx context.CancelFunc
}

type DashboardClient interface {
	FetchMessage() string
}

type DashboardNATSClient struct {
	conn *NATSConnection
}

func NewDashboardNATSClient() (*DashboardNATSClient, error) {
	natsConn, err := connectToNATS()
	if err != nil {
		return nil, err
	}
	return &DashboardNATSClient{conn: &natsConn}, nil
}

func (*DashboardNATSClient) FetchMessage() string {
	// have a NATSConnection with consumer object that can be stubbed for testing?
	natsConnection, err := connectToNATS()
	defer natsConnection.CancelCtx()
	defer natsConnection.NATSConn.Close()

	messageData := ""
	if err != nil {
		messageData += fmt.Sprintf("error during consumer initialization: %s\n", err)
		log.Print(messageData)
	}

	msgs, err := natsConnection.Consumer.Fetch(1)
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

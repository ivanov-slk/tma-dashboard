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

func FetchMessage() string {
	// have a server with consumer object that can be stubbed for testing?
	natsConnection, err := connectToNATS()
	if err != nil {
		log.Printf("error during consumer initialization: %s", err)
		defer natsConnection.CancelCtx()
		defer natsConnection.NATSConn.Close()
		return "hello message"
	}

	defer natsConnection.CancelCtx()
	defer natsConnection.NATSConn.Close()

	msgs, err := natsConnection.Consumer.Fetch(1)
	if err != nil {
		log.Fatalf("failed to fetch messages: %s", err)
	}

	messageData := "break-the-test-unless-a-message-is-received"
	for msg := range msgs.Messages() {
		msg.Ack()
		messageData = string(msg.Data())
	}
	if msgs.Error() != nil {
		fmt.Println("Error during Fetch(): ", msgs.Error())
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

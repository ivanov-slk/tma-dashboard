// Package natsclient provides functionalities for interacting with a NATS
// broker for fetching input data.
package natsclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func (d *DashboardNATSClient) Consume(out chan []byte) {
	cc, err := d.conn.Consumer.Consume(func(msg jetstream.Msg) {
		out <- msg.Data()
		msg.Ack()
	}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		fmt.Printf("failed to consume a message: %s", err)
	}))
	if err != nil {
		log.Fatalf("failed to start consuming: %s", err)
	}
	defer cc.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
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

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TMA",
		Subjects: []string{"generated-data"},
		Storage:  jetstream.MemoryStorage,
	})
	if err != nil {
		log.Fatalf("ERROR: failed to initialize stream: %s.", err)
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "TMADashboard",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("ERROR: failed to initialize consumer: %s.", err)
	}

	return NATSConnection{nc, c, ctx, cancel}
}

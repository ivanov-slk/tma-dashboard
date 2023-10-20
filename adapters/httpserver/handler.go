package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	nc, err := nats.Connect("http://nats-server:4222")
	if err != nil {
		log.Printf("failed to connect to nats: %s", err)
		fmt.Fprint(w, "hello message")
		return
	}
	defer nc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	js, _ := jetstream.New(nc)

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TMA",
		Subjects: []string{"generated-data"},
	})
	if err != nil {
		log.Fatalf("failed to get stream: %s", err)
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "TMA",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("failed to create consumer: %s", err)
	}

	msgs, err := c.Fetch(1)
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

	fmt.Fprint(w, messageData)
}

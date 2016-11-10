package pubsub

import (
	"context"
	"log"
)

var (
	broker *Broker
)

// Serve serves Nejireco Pub/Sub broker.
func Serve(ctx context.Context) {
	cfg := ConfigFromContext(ctx)

	log.Printf("Connecting redis server: %s", cfg.RedisURI)
	broker, err := NewBroker(ctx, cfg.RedisURI)
	if err != nil {
		log.Fatal(err)
	}
	defer broker.Close()

	err = broker.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = broker.SubscribeLocalTopics(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = broker.SubscribeCloudTopics(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = broker.ReceiveCloud(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = broker.ReceiveLocal(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

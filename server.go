package pubsub

import (
	"context"
	"log"
)

var (
	broker *Broker
)

// Serve serves Pub/Sub broker
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

	for _, t := range topics {
		broker.Subscribe(ctx, t)
	}
	broker.Receive(ctx)
}

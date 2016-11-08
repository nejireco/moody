package pubsub

import (
	"context"
	"log"
)

var (
	broker *Broker
)

func Serve(ctx context.Context) {
	cfg := ConfigFromContext(ctx)

	log.Printf("Connecting redis server: %s", cfg.RedisURI)
	broker, err := NewBroker(ctx, cfg.RedisURI)
	if err != nil {
		log.Fatal(err)
	}
	defer broker.Close()

	for _, t := range topics {
		broker.Subscribe(ctx, t)
	}
	broker.Receive(ctx)
}

package pubsub

import (
	"context"
	"log"
	"net/url"
	"strings"
	"sync"

	gcPubSub "cloud.google.com/go/pubsub"
	"github.com/garyburd/redigo/redis"
)

type Broker struct {
	PubConn        redis.Conn
	SubConn        redis.PubSubConn
	GCPubSubClient *gcPubSub.Client
	l              sync.RWMutex
}

func NewBroker(ctx context.Context, rawurl string) (*Broker, error) {
	pubConn, err := redis.DialURL(rawurl)
	if err != nil {
		return nil, err
	}

	subConn, err := redis.DialURL(rawurl)
	if err != nil {
		return nil, err
	}

	gcPubSubClient, err := NewGCPubSubClient(ctx)
	if err != nil {
		return nil, err
	}

	hc := &Broker{
		PubConn: pubConn,
		SubConn: redis.PubSubConn{
			Conn: subConn,
		},
		GCPubSubClient: gcPubSubClient,
	}

	return hc, nil
}

// Init initializes topics
func (b *Broker) Init(ctx context.Context) error {
	for _, topic := range topics {
		t, err := b.GCPubSubClient.CreateTopic(ctx, url.QueryEscape(topic))
		if err != nil {
			if strings.Contains(err.Error(), "Resource already exists in the project") {
				continue
			}
			return err
		}
		log.Printf("topic: %s created\n", t.ID())

	}
	return nil
}

// Subscribe subscribes topic
func (b *Broker) Subscribe(ctx context.Context, topic string) error {
	err := b.SubConn.Subscribe(topic)
	if err != nil {
		return err
	}
	return nil
}

// Receive emits pushed messages
func (b *Broker) Receive(ctx context.Context) {
	for {
		switch v := b.SubConn.Receive().(type) {
		case redis.Message:
			msgIDs, err := b.Emit(ctx, v.Channel, v.Data)
			if err != nil {
				log.Printf("Error: %s\n", err)
			} else {
				log.Printf("Emit topic: %s, data: %v, msgIDs: %v\n", v.Channel, v.Data, msgIDs)
			}
		case error:
			// TODO
			panic(v)
		}
	}
}

// Emit publishes all topics to Google Cloud Pub/Sub
func (b *Broker) Emit(ctx context.Context, topic string, data []byte) ([]string, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	t := b.GCPubSubClient.Topic(url.QueryEscape(topic))
	msgIDs, err := t.Publish(ctx, &gcPubSub.Message{
		Data: data,
	})
	if err != nil {
		// TODO
		return nil, err
	}
	return msgIDs, nil
}

// Close closes all connections
func (b *Broker) Close() error {
	if err := b.PubConn.Close(); err != nil {
		return err
	}
	if err := b.SubConn.Close(); err != nil {
		return err
	}
	if err := b.GCPubSubClient.Close(); err != nil {
		return err
	}
	return nil
}

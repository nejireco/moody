package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/garyburd/redigo/redis"
	"google.golang.org/api/iterator"
)

// Broker is a Nejireco Pub/Sub client.
type Broker struct {
	PubConn            redis.Conn
	SubConn            redis.PubSubConn
	CloudPubSubClient  *pubsub.Client
	CloudSubscriptions []*pubsub.Subscription
	l                  sync.RWMutex
}

// NewBroker creates a new broker client.
func NewBroker(ctx context.Context, rawurl string) (*Broker, error) {
	pubConn, err := redis.DialURL(rawurl)
	if err != nil {
		return nil, err
	}

	subConn, err := redis.DialURL(rawurl)
	if err != nil {
		return nil, err
	}

	cloudPubSubClient, err := NewCloudPubSubClient(ctx)
	if err != nil {
		return nil, err
	}

	hc := &Broker{
		PubConn: pubConn,
		SubConn: redis.PubSubConn{
			Conn: subConn,
		},
		CloudPubSubClient: cloudPubSubClient,
	}

	return hc, nil
}

// Init initializes topics.
func (b *Broker) Init(ctx context.Context) error {
	client := b.CloudPubSubClient
	cloudTopics = make(map[string]*pubsub.Topic)
	cloudSubscriptions = make(map[string]*pubsub.Subscription)
	for _, rawid := range topics {
		topic, err := CreateTopicIfNotExists(ctx, client, rawid)
		if err != nil {
			return err
		}
		cloudTopics[topic.ID()] = topic

		subscription, err := CreateSubscriptionIfNotExists(ctx, client, rawid, topic, 10*time.Second, nil)
		if err != nil {
			return err
		}
		cloudSubscriptions[subscription.ID()] = subscription
	}
	return nil
}

// SubscribeLocalTopics subscribes all topics published from local.
func (b *Broker) SubscribeLocalTopics(ctx context.Context) error {
	for id := range cloudTopics {
		rawid, err := url.QueryUnescape(id)
		if err != nil {
			return err
		}
		err = b.SubConn.Subscribe(rawid)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReceiveLocal receives pushed messages and publishes them to Google Cloud Pub/Sub.
func (b *Broker) ReceiveLocal(ctx context.Context) error {
	for {
		switch v := b.SubConn.Receive().(type) {
		case redis.Message:
			log.Printf("Receive from local msg: %s %s\n", v.Channel, string(v.Data))

			if MessageIsFromCloud(v.Data) {
				// Do not emit a message from Google Cloud Pub/Sub
				log.Println("Pass because from Cloud")
				continue
			}

			m := NewMessageFromLocal(v.Data)
			data, err := m.Marshal()
			if err != nil {
				log.Printf("Error on Marshal: %s", err)
				break
			}

			msgIDs, err := b.emitCloud(ctx, v.Channel, data)
			if err != nil {
				log.Printf("Error on EmitCloud: %s", err)
				return err
			}
			log.Printf("Emit to cloud topic: %s, data: %v, msgIDs: %v\n", v.Channel, v.Data, msgIDs)
		case error:
			return v
		}
	}
}

func (b *Broker) emitCloud(ctx context.Context, topic string, data []byte) ([]string, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	t, err := GetCloudTopic(ctx, topic)
	if err != nil {
		return nil, err
	}

	msgIDs, err := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("publish to cloud: %s, %s", t.ID(), data)
	return msgIDs, nil
}

// SubscribeCloudTopics subscribes all subscriptions published from Google Cloud Pub/Sub.
func (b *Broker) SubscribeCloudTopics(ctx context.Context) error {
	for id := range cloudSubscriptions {
		sub := cloudSubscriptions[id]
		b.CloudSubscriptions = append(b.CloudSubscriptions, sub)
	}
	return nil
}

// ReceiveCloud receives pushed messages and publishes them to local Pub/Sub.
func (b *Broker) ReceiveCloud(ctx context.Context) error {
	for _, sub := range b.CloudSubscriptions {
		go func(ctx context.Context, sub *pubsub.Subscription) {
			it, err := sub.Pull(ctx)
			if err != nil {
				log.Printf("Error on Pull: %s", err)
				return
			}
			defer it.Stop()

			for {
				msg, err := it.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Printf("Error on Next: %s", err)
					break
				}
				log.Printf("Receive from cloud msg: %s %s\n", sub.ID(), string(msg.Data))

				if MessageIsFromLocal(msg.Data) {
					// Do not emit a message from local Pub/Sub
					log.Println("Pass because from local")
					msg.Done(true)
					continue
				}

				m := NewMessageFromCloud(msg.Data)
				data, err := m.Marshal()
				if err != nil {
					log.Printf("Error on Marshal: %s", err)
					break
				}

				err = b.emitLocal(ctx, sub.ID(), data)
				if err != nil {
					log.Printf("Error on EmitLocal: %s", err)
					break
				}
				log.Printf("Emit to local topic: %s, data: %v\n", sub.ID(), msg.Data)
				msg.Done(true)
			}
		}(ctx, sub)
	}
	return nil
}

func (b *Broker) emitLocal(ctx context.Context, id string, data []byte) error {
	b.l.RLock()
	defer b.l.RUnlock()

	rawid, err := url.QueryUnescape(id)
	if err != nil {
		return err
	}
	_, err = b.PubConn.Do("PUBLISH", rawid, data)
	if err != nil {
		return err
	}
	log.Printf("publish to local: %s, %s", rawid, data)
	return nil
}

// Close closes all connections.
func (b *Broker) Close() error {
	if err := b.PubConn.Close(); err != nil {
		return err
	}
	if err := b.SubConn.Close(); err != nil {
		return err
	}
	if err := b.CloudPubSubClient.Close(); err != nil {
		return err
	}
	return nil
}

// Message wraps a Pub/Sub message.
type Message struct {
	Data []byte `json:"data"`
	From string `json:"__from__"`
}

// Marshal returns the JSON encoding of Message.
func (msg *Message) Marshal() ([]byte, error) {
	return json.Marshal(msg)
}

// NewMessageFromCloud creates a new Message from cloud.
func NewMessageFromCloud(data []byte) *Message {
	return &Message{
		Data: data,
		From: "cloud",
	}
}

// NewMessageFromLocal creates a new Message from local.
func NewMessageFromLocal(data []byte) *Message {
	return &Message{
		Data: data,
		From: "local",
	}
}

// MessageIsFromCloud reports whether bytes are Message from cloud.
func MessageIsFromCloud(b []byte) bool {
	m := &Message{}
	err := json.Unmarshal(b, m)
	return err == nil && m.From == "cloud"
}

// MessageIsFromLocal reports whether bytes are Message from local.
func MessageIsFromLocal(b []byte) bool {
	m := &Message{}
	err := json.Unmarshal(b, m)
	return err == nil && m.From == "local"
}

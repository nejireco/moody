package pubsub

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	gcPubSub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	cloudTopics map[string]*gcPubSub.Topic

	newCloudClient   = defaultNewCloudClientFunc
	createCloudTopic = defaultCreateCloudTopicFunc
	newCloudTopic    = defaultNewCloudTopicFunc
)

func defaultNewCloudClientFunc(ctx context.Context, projectID string, opts ...option.ClientOption) (*gcPubSub.Client, error) {
	return gcPubSub.NewClient(ctx, projectID, opts...)
}

func defaultCreateCloudTopicFunc(ctx context.Context, client *gcPubSub.Client, topic string) (*gcPubSub.Topic, error) {
	return client.CreateTopic(ctx, topic)
}

func defaultNewCloudTopicFunc(client *gcPubSub.Client, topic string) *gcPubSub.Topic {
	return client.Topic(topic)
}

// NewGCPubSubClient creates a new Google Cloud Pub/Sub client.
func NewGCPubSubClient(ctx context.Context) (*gcPubSub.Client, error) {
	cfg := ConfigFromContext(ctx)

	client, err := newCloudClient(ctx, cfg.GCP.ProjectID, option.WithServiceAccountFile(cfg.GCP.ServiceAccountFile))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateTopicIfNotExists creates a new Google Cloud Pub/Sub topic if it doesn't exist.
func CreateTopicIfNotExists(ctx context.Context, client *gcPubSub.Client, rawtopic string) (*gcPubSub.Topic, error) {
	topic := url.QueryEscape(rawtopic)
	t, err := createCloudTopic(ctx, client, topic)
	if err != nil {
		if strings.Contains(err.Error(), "Resource already exists in the project") {
			return newCloudTopic(client, topic), nil
		}
		return nil, err
	}
	return t, nil
}

// GetCloudTopic returns registered topic.
func GetCloudTopic(ctx context.Context, topic string) (*gcPubSub.Topic, error) {
	if t, ok := cloudTopics[url.QueryEscape(topic)]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("topic %s not found", topic)
}

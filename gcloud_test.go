package pubsub

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"

	gcPubSub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func TestNewGCPubSubClient(t *testing.T) {
	ctx := context.Background()
	projectID := "test-project"
	saFile := "./sa.json"
	cfg := &Config{
		GCP: &GCPConfig{
			ProjectID:          projectID,
			ServiceAccountFile: saFile,
		},
	}
	ctx = NewContext(ctx, cfg)
	orgFunc := newCloudClient
	defer func() {
		newCloudClient = orgFunc
	}()

	var argProjectID string
	var argOpt option.ClientOption
	newCloudClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (*gcPubSub.Client, error) {
		argProjectID = projectID
		argOpt = opts[0]
		return nil, nil
	}
	_, err := NewGCPubSubClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if argProjectID != projectID {
		t.Errorf("expected: %s, but got: %s", projectID, argProjectID)
	}
	if fmt.Sprint(argOpt) != saFile {
		t.Errorf("expected: %s, but got: %s", saFile, argOpt)
	}
}

func TestCreateTopicIfNotExists(t *testing.T) {
	ctx := context.Background()
	orgFunc1 := createCloudTopic
	orgFunc2 := newCloudTopic
	defer func() {
		createCloudTopic = orgFunc1
		newCloudTopic = orgFunc2
	}()

	var argTopic string
	rawtopic := "test/topic"
	topic := url.QueryEscape(rawtopic)
	createCloudTopic = func(ctx context.Context, client *gcPubSub.Client, topic string) (*gcPubSub.Topic, error) {
		argTopic = topic
		return &gcPubSub.Topic{}, nil
	}
	_, err := CreateTopicIfNotExists(ctx, nil, rawtopic)
	if err != nil {
		t.Fatal(err)
	}
	if argTopic != topic {
		t.Errorf("expected: %s, but got: %s", topic, argTopic)
	}

	createCloudTopic = func(ctx context.Context, client *gcPubSub.Client, topic string) (*gcPubSub.Topic, error) {
		return nil, errors.New("Resource already exists in the project (resource=new-topic)")
	}
	newCloudTopic = func(client *gcPubSub.Client, topic string) *gcPubSub.Topic {
		argTopic = topic
		return nil
	}
	_, err = CreateTopicIfNotExists(ctx, nil, rawtopic)
	if err != nil {
		t.Fatal(err)
	}
	if argTopic != topic {
		t.Errorf("expected: %s, but got: %s", topic, argTopic)
	}
}

func TestGetCloudTopic(t *testing.T) {
	ctx := context.Background()
	cloudTopics = make(map[string]*gcPubSub.Topic)
	rawtopic := "test/topic"
	topic := url.QueryEscape(rawtopic)
	cloudTopics[topic] = &gcPubSub.Topic{}

	got, err := GetCloudTopic(ctx, rawtopic)
	if err != nil {
		t.Errorf("expected: nil, but got: %v", err)
	}
	if got == nil {
		t.Errorf("expected: not nil, but got: %#v", got)
	}

	got, err = GetCloudTopic(ctx, topic)
	if err == nil {
		t.Errorf("expected: not nil, but got: %v", err)
	}
	if got != nil {
		t.Errorf("expected: nil, but got: %#v", got)
	}
}

package pubsub

import (
	"context"

	gcPubSub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// NewGCPubSubClient creates a new Google Cloud Pub/Sub client.
func NewGCPubSubClient(ctx context.Context) (*gcPubSub.Client, error) {
	cfg := ConfigFromContext(ctx)

	client, err := gcPubSub.NewClient(ctx, cfg.GCP.ProjectID, option.WithServiceAccountFile(cfg.GCP.ServiceAccountFile))
	if err != nil {
		return nil, err
	}

	return client, nil
}

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"

	pb "github.com/kolya59/easy_normalization/proto"
)

type Client struct {
	topic *pubsub.Topic
}

func NewClient(projectID, topicName string) (*Client, error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	topic := client.Topic(topicName)

	// Create the topic if it doesn't exist
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check topic existense: %v", err)
	}
	if !exists {
		if _, err = client.CreateTopic(ctx, topicName); err != nil {
			return nil, fmt.Errorf("failed to create topic")
		}
	}

	return &Client{topic: topic}, nil
}

func (c *Client) SendCars(cars []pb.Car) error {
	data, err := json.Marshal(cars)
	if err != nil {
		return err
	}
	msg := &pubsub.Message{Data: data}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := c.topic.Publish(ctx, msg).Get(ctx); err != nil {
		return err
	}

	return nil
}

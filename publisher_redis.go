package gqlgen_subscription_redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var _ Publisher = (*redisPublisher)(nil)

type redisPublisher struct {
	client       *redis.Client
	redisChannel string
}

// NewRedisPublisher new publisher using redis
func NewRedisPublisher(client *redis.Client, redisChannel string) Publisher {
	return &redisPublisher{
		client:       client,
		redisChannel: redisChannel,
	}
}

func (p redisPublisher) Publish(channelName string, data []byte) error {
	payload := Payload{
		ChannelName: channelName,
		Data:        data,
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.client.Publish(context.TODO(), p.redisChannel, bytes).Err()
}

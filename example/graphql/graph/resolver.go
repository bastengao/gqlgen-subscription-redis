package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"fmt"

	sub "github.com/bastengao/gqlgen-subscription-redis"
	"github.com/go-redis/redis/v8"
)

type Resolver struct {
	broker    sub.Broker
	publisher sub.Publisher
}

func NewResolver() *Resolver {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	redisChannel := "graphql_subscription_channel"

	broker := sub.NewRedisBroker(client, redisChannel)
	publisher := sub.NewRedisPublisher(client, redisChannel)

	go func() {
		fmt.Println(broker.Receive())
	}()
	return &Resolver{
		broker:    broker,
		publisher: publisher,
	}
}

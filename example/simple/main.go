package main

import (
	"fmt"

	sub "github.com/bastengao/gqlgen-subscription-redis"

	"github.com/go-redis/redis/v8"
)

const redisChannel = "graphql_subscription_channel"
const graphqlChannel = "chat_room"

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	var subscriptionBroker = sub.NewRedisBroker(client, redisChannel)
	messageCh := make(chan string)
	unsub, err := subscriptionBroker.Subscribe(graphqlChannel, messageCh, func(id string, payload sub.Payload) (interface{}, error) {
		return string(payload.Data), nil
	})
	if err != nil {
		panic(err)
	}
	defer unsub.Close()

	go func() {
		fmt.Println(subscriptionBroker.Receive())
	}()

	var publisher sub.Publisher = sub.NewRedisPublisher(client, redisChannel)
	err = publisher.Publish(graphqlChannel, []byte("new message"))
	if err != nil {
		panic(err)
	}

	fmt.Println(<-messageCh)
}

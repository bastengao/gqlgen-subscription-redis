package gqlgen_subscription_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/go-redis/redis/v8"
)

var _ Broker = (*redisBroker)(nil)

type redisBroker struct {
	client        *redis.Client
	redisChannel  string
	subscriptions sync.Map // channelName => channelSubscriptions
}

// NewRedisBroker new broker using redis
func NewRedisBroker(client *redis.Client, redisChannel string) Broker {
	return &redisBroker{
		client:       client,
		redisChannel: redisChannel,
	}
}

func (b *redisBroker) Receive() error {
	pubsub := b.client.Subscribe(context.TODO(), b.redisChannel)

	_, err := pubsub.Receive(context.TODO())
	if err != nil {
		return err
	}

	redisChannel := pubsub.Channel()
	for msg := range redisChannel {
		var payload Payload
		err := json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			return err
		}

		chValue, exist := b.subscriptions.Load(payload.ChannelName)
		if !exist {
			continue
		}
		ch, ok := chValue.(*channelSubscriptions)
		if !ok {
			continue
		}

		ch.RLock()
		for id, subscription := range ch.subscriptions {
			message, err := subscription.handler(id, payload)
			if err != nil {
				log.Println(err)
				continue
			}

			subscription.channel.Send(reflect.ValueOf(message))
		}
		ch.RUnlock()
	}

	return pubsub.Close()
}

func (b *redisBroker) Subscribe(channelName string, channel interface{}, handler PayloadHandler) (Unsubscriber, error) {
	value := reflect.ValueOf(channel)
	if kind := value.Kind(); kind != reflect.Chan {
		return nil, fmt.Errorf("channel must be channel, but %v", kind)
	}

	chValue, ok := b.subscriptions.Load(channelName)
	if !ok {
		newChannel := &channelSubscriptions{subscriptions: make(map[string]*subscription)}
		chValue, _ = b.subscriptions.LoadOrStore(channelName, newChannel)
	}

	ch, _ := chValue.(*channelSubscriptions)
	return ch.append(channelName, value, handler), nil
}

type subscription struct {
	channelName string
	channel     reflect.Value
	handler     PayloadHandler
}

type channelSubscriptions struct {
	sync.RWMutex
	subscriptions map[string]*subscription // subscription id => subscription
}

func (c *channelSubscriptions) append(channelName string, channel reflect.Value, handler PayloadHandler) Unsubscriber {
	c.Lock()
	defer c.Unlock()
	var id string
	for {
		id = randomNumberStr(6)
		_, exists := c.subscriptions[id]
		if !exists { // no conflict
			break
		}
	}

	c.subscriptions[id] = &subscription{
		channelName: channelName,
		channel:     channel,
		handler:     handler,
	}

	return CloserFunc(func() error {
		c.remove(id)
		return nil
	})
}

func (c *channelSubscriptions) remove(id string) {
	c.Lock()
	defer c.Unlock()
	delete(c.subscriptions, id)
}

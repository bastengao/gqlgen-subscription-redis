package gqlgen_subscription_redis

import (
	"context"
	"io"
	"log"
	"reflect"
)

// Payload represents payload sending to subscription
type Payload struct {
	ChannelName string
	Data        []byte
}

// PayloadHandler represents function that returns message to subscription channel from payload
type PayloadHandler func(subscriptionID string, payload Payload) (interface{}, error)

// Unsubscriber is the interface the wrap the Close method to unsubscribe
type Unsubscriber io.Closer

// CloserFunc is an adapter to allow use of ordinary functions as io.Closer
type CloserFunc func() error

// Close calls c
func (c CloserFunc) Close() error {
	return c()
}

// CloseSubscription unsubscribes when subscription is closed
func CloseSubscription(ctx context.Context, u Unsubscriber, channel interface{}) {
	<-ctx.Done()
	u.Close()

	cValue := reflect.ValueOf(channel)
	if cValue.Kind() == reflect.Chan {
		cValue.Close()
	} else {
		log.Printf("type of channel is not channel but %s", cValue.Kind())
	}
}

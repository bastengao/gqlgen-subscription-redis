package gqlgen_subscription_redis

// Publisher ...
type Publisher interface {
	// Publish publishes payload data to channel
	Publish(channelName string, data []byte) error
}

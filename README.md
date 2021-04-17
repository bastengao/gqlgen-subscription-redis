# gqlgen-subscription-redis

Use redis pub/sub to send graphql subscription message.

## Example

Full [example](https://github.com/bastengao/gqlgen-subscription-redis/tree/master/example/graphql).


### Subscribe

```go
func (r *subscriptionResolver) NewMessagePosted(ctx context.Context, chatRoom string) (<-chan string, error) {
	ch := make(chan string)
	unsub, err := r.broker.Subscribe(chatRoom, ch, func(id string, payload sub.Payload) (interface{}, error) {
		return string(payload.Data), nil
	})
	if err != nil {
		return nil, err
	}

	go sub.CloseSubscription(ctx, unsub, ch)
	return ch, nil
}
```

### Publish

```go
func (r *mutationResolver) PostMessage(ctx context.Context, chatRoom string, message string) (bool, error) {
	err := r.publisher.Publish(chatRoom, []byte(message))
	if err != nil {
		return false, err
	}
	return true, nil
}
```

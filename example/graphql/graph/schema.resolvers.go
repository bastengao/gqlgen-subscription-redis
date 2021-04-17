package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"example/graphql/graph/generated"

	sub "github.com/bastengao/gqlgen-subscription-redis"
)

func (r *mutationResolver) PostMessage(ctx context.Context, chatRoom string, message string) (bool, error) {
	err := r.publisher.Publish(chatRoom, []byte(message))
	if err != nil {
		return false, err
	}
	return true, nil
}

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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

package main

import (
	"context"

	"github.com/ml444/gomega/broker"
)

func pub(ctx context.Context, topic string, item *broker.Item) error {
	// file sequence offset
	bk := broker.GetOrCreateBroker(topic)
	err := bk.Send(item)
	if err != nil {
		return err
	}
	return nil
}

package utils

import (
	"context"

	"cloud.google.com/go/storage"
)

func NewObjectStorage() (*storage.Client, error) {

	var ctx context.Context
	var client *storage.Client
	var err error

	ctx = context.Background()

	client, err = storage.NewClient(ctx)

	return client, err

}

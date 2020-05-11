package storage

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
)

var (
	projectID       = os.Getenv("GOOGLE_CLOUD_PROJECT")
	k_service       = os.Getenv("K_SERVICE")
	k_revision      = os.Getenv("K_REVISION")
	k_configuration = os.Getenv("K_CONFIGURATION")
)

func StorageClient(ctx context.Context) *storage.Client {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return client
}

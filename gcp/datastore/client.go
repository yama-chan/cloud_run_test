package datastore

import (
	"context"
	"log"
	"os"

	// "github.com/labstack/echo/middleware"

	// "github.com/labstack/echo/v4"

	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
	"cloud.google.com/go/datastore"
	// "github.com/labstack/echo"
)

var (
	projectID       = os.Getenv("GOOGLE_CLOUD_PROJECT")
	k_service       = os.Getenv("K_SERVICE")
	k_revision      = os.Getenv("K_REVISION")
	k_configuration = os.Getenv("K_CONFIGURATION")
)

type Datastore struct {
	client *datastore.Client
}

func NewDatastore(ctx context.Context) Datastore {
	return Datastore{
		client: datastoreClient(ctx),
	}
}

func datastoreClient(ctx context.Context) *datastore.Client {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalln("fail to datastoreClient :" + err.Error())
	}
	return client
}

package firestore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var (
	projectID       = os.Getenv("GOOGLE_CLOUD_PROJECT")
	k_service       = os.Getenv("K_SERVICE")
	k_revision      = os.Getenv("K_REVISION")
	k_configuration = os.Getenv("K_CONFIGURATION")
)

type Firestore struct {
	client *firestore.Client
}

func NewFirestore(ctx context.Context) Firestore {
	return Firestore{
		client: firestoreClient(ctx),
	}
}

func firestoreClient(ctx context.Context) *firestore.Client {
	// Use the application default credentials
	// ctx := context.Background()
	conf := &firebase.Config{
		ProjectID: projectID,
	}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// defer client.Close()
	return client
}

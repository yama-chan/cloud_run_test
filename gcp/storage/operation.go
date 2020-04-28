package firestore

import (
	"log"

	"cloud.google.com/go/firestore"
)

func (*firestore.Client) Insert() {
	// データ追加
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

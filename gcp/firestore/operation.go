package firestore

import (
	"context"
	"log"
)

func (f Firestore) Insert(ctx context.Context) error {
	// データ追加
	_, _, err := f.client.Collection("test").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	return err
}

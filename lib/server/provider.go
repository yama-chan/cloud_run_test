package server

import (
	"context"
	"net/http"

	"github.com/taisukeyamashita/test/gcp/datastore"
	"github.com/taisukeyamashita/test/gcp/storage"
	testusecase "github.com/taisukeyamashita/test/usecase/test"
)

// Provider アプリケーションで使用するユースケースの提供を行う関数をインターフェースで定義
type Provider interface {
	TestUsecase(ctx context.Context) *testusecase.TestUsecase
	// TestUsecase1(ctx context.Context)
	// TestUsecase2(ctx context.Context)
	Context(r *http.Request) context.Context
	Finalize(ctx context.Context) error
	ProvideStorageOperator(ctx context.Context) storage.StorageOpeator
	ProvideDatastoreOperator(ctx context.Context) datastore.DatastoreOperator
}

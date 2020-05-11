package server

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	datastore2 "github.com/taisukeyamashita/test/gcp/datastore"
	storage2 "github.com/taisukeyamashita/test/gcp/storage"
	"github.com/taisukeyamashita/test/lib/errs"
	"github.com/taisukeyamashita/test/lib/server"
	testusecase "github.com/taisukeyamashita/test/usecase/test"
)

type (
	datastoreClientKeyString   string
	storageClientKeyString     string
	contextFinalizersKeyString string
)

var (
	datastoreClientKey   datastoreClientKeyString   = "datastoreClientKey"
	storageClientKey     storageClientKeyString     = "storageClient"
	contextFinalizersKey contextFinalizersKeyString = "context_finalizers"
)

// AppProvider ユースケースの提供を行うプロバイダー implement Provider interface
type AppProvider struct{}

var _ server.Provider = &AppProvider{}

type finalizer func() error

// NewAppProvider アプリケーションのプロバイダーを新規生成
func NewAppProvider() server.Provider {
	return &AppProvider{}
}

// TestUsecase テストユースケースを提供
func (p *AppProvider) TestUsecase(ctx context.Context) *testusecase.TestUsecase {
	return testusecase.ProvideTestUsecase(
		p.provideStorageOperator(ctx),
		p.provideDatastoreOperator(ctx),
	)
}

// Context implements server.Provider
// http.Requestを受けてアプリケーション用のcontext.Contextの成形して返すようにする
func (p *AppProvider) Context(r http.Request) context.Context {
	ctx := r.Context()
	// context からGCPサービスのclientとfinalizer関数を生成
	// しかし、すでにcontext.Contextにclientがあるならば新規で作成しないように実装すること
	var (
		datastore, dsFinalizer    = p.datastoreClientWithFinalizer(ctx)
		storage, storageFinalizer = p.storageClientWithFinalizer(ctx)
	)
	// context に各clientを格納
	ctx = context.WithValue(ctx, datastoreClientKey, datastore)
	ctx = context.WithValue(ctx, storageClientKey, storage)
	// context に各clientのfinalizer関数をリスト形式で格納
	ctx = context.WithValue(ctx, contextFinalizersKey,
		[]finalizer{
			dsFinalizer,
			storageFinalizer,
		})
	return ctx
}

// Finalize implements server.Provider
func (p *AppProvider) Finalize(ctx context.Context) error {
	// context に格納したfinalizer関数を取得
	finalizers, ok := ctx.Value(contextFinalizersKey).([]finalizer)
	if !ok {
		return nil
	}
	errMsgs := make([]string, 0)
	for _, f := range finalizers {
		if err := f(); err != nil {
			errMsgs = append(errMsgs, err.Error())
		}
	}
	if len(errMsgs) != 0 {
		return errs.NewXerror(strings.Join(errMsgs, ";"))
	}
	return nil
}

func (p *AppProvider) provideStorageOperator(ctx context.Context) storage2.StorageOpeator {
	client, _ := p.storageClientWithFinalizer(ctx)
	return storage2.ProvideStorageOpeator(client)
}

func (p *AppProvider) storageClientWithFinalizer(ctx context.Context) (*storage.Client, finalizer) {
	closeFn := func(c *storage.Client) finalizer { // finalizer func() error関数を返す関数をスコープ内に生成
		return func() error {
			go c.Close()
			return nil
		}
	}
	//すでにcontext.Contextに格納されているclientがあるならば新規で作成しない
	//context.Contextに格納している*storage.Clientを取得
	client, ok := ctx.Value(storageClientKey).(*storage.Client)
	if ok {
		return client, closeFn(client)
	}
	//Context(r http.Request)が呼ばれなければここが通る
	//context.Contextに格納していなければ*storage.Clientを新規作成
	client = storage2.StorageClient(ctx)
	return client, closeFn(client)
}

func (p *AppProvider) provideDatastoreOperator(ctx context.Context) *datastore2.Operator {
	client, _ := p.datastoreClientWithFinalizer(ctx)
	return datastore2.ProvideDatastoreOperator(client)
}

func (p *AppProvider) datastoreClientWithFinalizer(ctx context.Context) (*datastore.Client, finalizer) {
	closeFn := func(c *datastore.Client) finalizer { // finalizer func() error関数を返す関数をスコープ内に生成
		return func() error {
			go c.Close()
			return nil
		}
	}
	//すでにcontext.Contextに格納されているclientがあるならば新規で作成しない
	//context.Contextに格納している*storage.Clientを取得
	client, ok := ctx.Value(datastoreClientKey).(*datastore.Client)
	if ok {
		return client, closeFn(client)
	}
	//Context(r http.Request)が呼ばれなければここが通る
	//context.Contextに格納していなければ*storage.Clientを新規作成
	client = datastore2.DatastoreClient(ctx)
	return client, closeFn(client)
}

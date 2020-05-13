package provider

import (
	"context"
	"net/http"
	"strings"

	datastore3 "cloud.google.com/go/datastore"
	storage3 "cloud.google.com/go/storage"
	"github.com/taisukeyamashita/test/gcp/datastore"
	"github.com/taisukeyamashita/test/gcp/storage"
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
type AppProvider struct {
	Gcs       storage.StorageOpeator
	Datastore datastore.DatastoreOperator
}

type finalizer func() error

// NewAppProvider アプリケーションのプロバイダーを新規生成
func NewAppProvider() *AppProvider {
	return &AppProvider{}
}

// AppProvider implements server.Provider
var _ server.Provider = &AppProvider{}

// TestUsecase テストユースケースを提供
func (p *AppProvider) TestUsecase(ctx context.Context) *testusecase.TestUsecase {
	return testusecase.ProvideTestUsecase(ctx, p.ProvideStorageOperator(ctx), p.ProvideDatastoreOperator(ctx))
}

// Context implements server.Provider
// http.Requestを受けてアプリケーション用のcontext.Contextの成形して返すようにする
func (p *AppProvider) Context(r *http.Request) context.Context {
	ctx := r.Context()
	// context からGCPサービスのclientとfinalizer関数を生成
	// しかし、現状の実装ならcontext.Contextにclientがあるならば新規で作成しないように実装することになる
	var (
		datastore, dsFinalizer    = p.datastoreClientWithFinalizer(ctx)
		storage, storageFinalizer = p.storageClientWithFinalizer(ctx)
	)
	// TODO:　 clientはできればProviderに持たせるように修正する
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
		return errs.NewXerrorWithMessage(strings.Join(errMsgs, ";"))
	}
	return nil
}

func (p *AppProvider) ProvideStorageOperator(ctx context.Context) storage.StorageOpeator {
	client, _ := p.storageClientWithFinalizer(ctx)
	return storage.ProvideStorageOpeator(client)
}

func (p *AppProvider) storageClientWithFinalizer(ctx context.Context) (*storage3.Client, finalizer) {
	closeFn := func(c *storage3.Client) finalizer { // finalizer func() error関数を返す関数をスコープ内に生成
		return func() error {
			go c.Close()
			return nil
		}
	}
	//すでにcontext.Contextに格納されているclientがあるならば新規で作成しない
	//context.Contextに格納している*storage3.Clientを取得
	client, ok := ctx.Value(storageClientKey).(*storage3.Client)
	if ok {
		return client, closeFn(client)
	}
	//Context(r http.Request)が呼ばれなければここが通る
	//context.Contextに格納していなければ*storage3.Clientを新規作成
	client = storage.StorageClient(ctx)
	return client, closeFn(client)
}

func (p *AppProvider) ProvideDatastoreOperator(ctx context.Context) datastore.DatastoreOperator {
	client, _ := p.datastoreClientWithFinalizer(ctx)
	return datastore.ProvideDatastoreOperator(client)
}

func (p *AppProvider) datastoreClientWithFinalizer(ctx context.Context) (*datastore3.Client, finalizer) {
	closeFn := func(c *datastore3.Client) finalizer { // finalizer func() error関数を返す関数をスコープ内に生成
		return func() error {
			go c.Close()
			return nil
		}
	}
	//すでにcontext.Contextに格納されているclientがあるならば新規で作成しない
	//context.Contextに格納している*storage3.Clientを取得
	client, ok := ctx.Value(datastoreClientKey).(*datastore3.Client)
	if ok {
		return client, closeFn(client)
	}
	//Context(r http.Request)が呼ばれなければここが通る
	//context.Contextに格納していなければ*storage3.Clientを新規作成
	client = datastore.DatastoreClient(ctx)
	return client, closeFn(client)
}

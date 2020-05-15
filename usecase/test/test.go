package testusecase

import (
	"context"

	"github.com/taisukeyamashita/test/gcp/datastore"
	"github.com/taisukeyamashita/test/gcp/storage"
	"github.com/taisukeyamashita/test/lib/errs"
	"github.com/taisukeyamashita/test/lib/times"
	utils "github.com/taisukeyamashita/test/lib/util"
	"github.com/taisukeyamashita/test/lib/validate"
	"github.com/taisukeyamashita/test/model"
	"github.com/taisukeyamashita/test/usecase"
)

type TestUsecase struct {
	StorageOperator   storage.StorageOpeator
	DatastoreOperator datastore.DatastoreOperator
}

var _ usecase.Usecase = &TestUsecase{}

func ProvideTestUsecase(ctx context.Context, s storage.StorageOpeator, d datastore.DatastoreOperator) *TestUsecase {
	return &TestUsecase{
		StorageOperator:   s,
		DatastoreOperator: d,
	}
}

// テスト関連のユースケースを以下に用意する

// Insert ユーザデータを格納するユースケース
func (t TestUsecase) Insert(ctx context.Context) error {
	const format = "2006/01/02 15:04:05" // 24h表現、0埋めあり
	user := model.UserInf{
		ID:               utils.CreateUniqueId(),
		Fullname:         "test２",
		LastModifiedDate: times.CurrentTime().Format(format),
	}
	validate.ValidateStruct(user) // バリデート
	putErr := t.DatastoreOperator.Put(ctx, &user)
	if putErr != nil {
		//  「: %v」 とすることで既存のerrorの情報を出力する
		return errs.WrapXerror(putErr)
	}
	return nil
}

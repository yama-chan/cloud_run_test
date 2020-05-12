package datastore

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/taisukeyamashita/test/lib/config"
	"github.com/taisukeyamashita/test/lib/errs"
	utils "github.com/taisukeyamashita/test/lib/util"
	"github.com/taisukeyamashita/test/model"
)

type DatastoreOperator interface {
	DatastoreClient() *datastore.Client
	Put(ctx context.Context, userInf *model.UserInf) error
}

type Config struct {
	ProjectID string
}

type Operator struct {
	Client *datastore.Client
	Config Config
}

var _ DatastoreOperator = &Operator{}

func ProvideDatastoreOperator(client *datastore.Client) *Operator {
	return &Operator{
		Client: client,
		Config: Config{
			ProjectID: config.ProjectID,
		},
	}
}

//lessonのスライスを型宣言
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type Lessons []Lesson

var lessonInf *Lesson

//lessonの構造体
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type Lesson struct {
	mu               sync.Mutex     //排他制御用　※一意のIDを生成するためユーザ同時登録を防ぐ
	key              *datastore.Key //データのkey ※更新時に使用
	ID               string
	LessonName       string
	Communication    string
	Analytics        string
	Security         string
	Description      string
	RegisteredDate   string
	LastModifiedDate string
}

// type TestUser datastore.Entity{
// 	kind string
// }

func (op *Operator) DatastoreClient() *datastore.Client {
	return op.Client
}

func (op *Operator) Put(ctx context.Context, userInf *model.UserInf) error {
	var (
		taskKey *datastore.Key
		err     error
	)
	kind := "user"
	if userInf.EncodedKey != "" {
		enKey := userInf.EncodedKey
		taskKey, err = datastore.DecodeKey(enKey)
		if err != nil {
			s := []string{"Failed to DecodeKey:", err.Error()}
			log.Println(strings.Join(s, ""))
			return errs.WrapXerror(err)
		}
	} else {
		// taskKey = datastore.IncompleteKey(kind, nil)
		taskKey = datastore.NameKey(kind, utils.CreateUniqueId(), nil)
	}
	//保存対象の構造体のポインタを引数に定義しPUT
	if taskKey != nil {
		_, err := op.Client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			//keyを暗号化して保存（更新用）
			userInf.EncodedKey = taskKey.Encode()
			if _, err := tx.Put(taskKey, userInf); err != nil {
				return err
			}
			fmt.Printf("Saved %v: %v\n", taskKey, userInf.Fullname)
			return nil
		})
		if err != nil {
			return errs.WrapXerror(err)
		}
	}
	return nil
}

func (ds Datastore) Put(ctx context.Context, userInf *model.UserInf) error {
	var (
		taskKey *datastore.Key
		err     error
	)
	// Sets the kind for the new entity.
	kind := "user"
	// // Sets the name/ID for the new entity.
	// name := ""

	// keyインスタンスの生成
	// taskKey := datastore.NameKey(kind, name, nil)
	if userInf.EncodedKey != "" {
		enKey := userInf.EncodedKey
		taskKey, err = datastore.DecodeKey(enKey)
		if err != nil {
			s := []string{"Failed to DecodeKey:", err.Error()}
			log.Println(strings.Join(s, ""))
			return err
		}
	} else {
		taskKey = datastore.IncompleteKey(kind, nil)
	}
	fmt.Printf("taskKey >>>>>>>> %v\n", taskKey)
	fmt.Printf("client >>>>>>>> %v\n", ds.client)

	//保存対象の構造体のポインタを引数に定義しPUT
	if taskKey != nil {
		//keyを暗号化して保存（更新用）
		userInf.EncodedKey = taskKey.Encode()
		if _, err := ds.client.Put(ctx, taskKey, userInf); err != nil {
			s := []string{"Failed to save task:", err.Error()}
			log.Println(strings.Join(s, ""))
			return err
		}
		fmt.Printf("Saved %v: %v\n", taskKey, userInf.Fullname)
	}
	return nil
}

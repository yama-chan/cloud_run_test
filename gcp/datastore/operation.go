package datastore

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/taisukeyamashita/test/model"
)

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

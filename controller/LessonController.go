package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	// "github.com/labstack/echo/middleware"

	// "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4"

	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
	"cloud.google.com/go/datastore"
	// "github.com/labstack/echo"
)

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

//lessonのスライスを型宣言
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type Lessons []Lesson

var lessonInf *Lesson

func AddLesson(c echo.Context) error {

	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	projectID := "sandbox-taisukeyamashita"

	log.Printf("ctx >>>>>>>>>>>>>>>>>>>>>>> %v", ctx)

	r := c.Request()
	r.ParseForm() //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。

	//日付フォーマット
	const format = "2006/01/02 15:04:05" // 24h表現、0埋めあり

	//IDフォーマット
	const formatForID = "20060102150405"

	lessonId := strings.Join(r.Form["lessonId"], "")
	log.Printf("userId >>>>>>>>>>>>>>>>>>>>>>> %v\n", lessonId)

	//更新の場合
	if lessonId != "" {
		log.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 更新\n")
		lesson, err := GetLessonByID(c, lessonId)
		lessonInf = &lesson

		if err != nil {
			c.Render(http.StatusOK, "error", err.Error())
			return err
		}
		//新規登録の場合
	} else {
		log.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 新規\n")
		//保存対象となる構造体を定義,　ポイント型（&user）としないこと
		//《注意》：PUTの引数は構造体のポインタ(&user)とすること
		lesson := Lesson{}
		lessonInf = &lesson
	}

	//保存対象となる構造体を定義,　ポイント型（&user）としないこと
	//《注意》：PUTの引数は構造体のポインタ(&user)とすること
	// lesson := Lesson{
	// 	LessonName:     strings.Join(r.Form["lessonname"], ""),
	// 	Communication:  strings.Join(r.Form["communication"], ""),
	// 	Analytics:      strings.Join(r.Form["analytics"], ""),
	// 	Security:       strings.Join(r.Form["security"], ""),
	// 	Description:    strings.Join(r.Form["comment"], ""),
	// 	RegisteredDate: time.Now().Format(format)}

	//一意のID生成時処理のため排他ロック開始
	lessonInf.mu.Lock()

	//関数がreturnするまではロック
	defer lessonInf.mu.Unlock()

	//登録日時秒の組み合わせは一意とするためロック内で定義
	registeredDate := time.Now()

	lessonInf.LessonName = strings.Join(r.Form["lessonname"], "")
	lessonInf.Communication = strings.Join(r.Form["communication"], "")
	lessonInf.Analytics = strings.Join(r.Form["analytics"], "")
	lessonInf.Security = strings.Join(r.Form["security"], "")
	lessonInf.Description = strings.Join(r.Form["comment"], "")
	// lessonInf.RegisteredDate = registeredDate.Format(format)

	//IDの無いユーザはIDを割り振って登録日時を更新する。
	if lessonInf.ID == "" {
		//ID用の文字列スライスを生成
		s := []string{"L-", registeredDate.Format(formatForID)}

		lessonInf.RegisteredDate = registeredDate.Format(format)
		lessonInf.ID = strings.Join(s, "")
		//すでにIDが割り振られているユーザは最終更新日時を変更する。
	} else {
		lessonInf.LastModifiedDate = registeredDate.Format(format)
	}

	fmt.Println(r.Form) //これらのデータはサーバのプリント情報に出力されます
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	//r.Form["fullname"]ではvalueを配列[]で取得するため size=1 でもrangeまたは配列要素を結合するなどの処理が必要となる
	fmt.Println(r.Form["lessonname"])
	fmt.Println(r.Form["communication"])
	fmt.Println(r.Form["analytics"])
	fmt.Println(r.Form["security"])
	fmt.Println(r.Form["comment"])

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return err
	}

	// Sets the kind for the new entity.
	kind := "lesson"
	// // Sets the name/ID for the new entity.
	// name := ""

	// keyインスタンスの生成
	// taskKey := datastore.NameKey(kind, name, nil)
	taskKey := datastore.IncompleteKey(kind, nil)

	fmt.Printf("taskKey >>>>>>>> %v\n", taskKey)
	fmt.Printf("client >>>>>>>> %v\n", client)

	//保存対象の構造体のポインタを引数に定義しPUT
	if lessonInf.key != nil { //更新の場合　※エンティティがすでにkeyを持っている
		if _, err := client.Put(ctx, lessonInf.key, lessonInf); err != nil {
			s := []string{"Failed to save task:", err.Error()}
			c.Render(http.StatusOK, "error", strings.Join(s, ""))
			return err
		}

	} else { //新規の場合　※新規のkey(taskKey)を生成し保村処理を行う必要がある
		if _, err := client.Put(ctx, taskKey, lessonInf); err != nil {
			s := []string{"Failed to save task:", err.Error()}
			c.Render(http.StatusOK, "error", strings.Join(s, ""))
			return err
		}
	}

	fmt.Printf("Saved %v: %v\n", taskKey, lessonInf.LessonName)

	return nil
}

func GetLesson(c echo.Context) []Lesson {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		// TODO: Handle error.
	}

	var lessons Lessons

	// query := datastore.NewQuery("lesson").Filter("IsDeleted =", false)
	query := datastore.NewQuery("lesson")

	keys, err := client.GetAll(ctx, query, &lessons)
	if err != nil {
		// TODO: Handle error.
	}
	for i, key := range keys {
		fmt.Println(key)
		fmt.Println(lessons[i])
	}

	return lessons
}

func GetLessonByID(c echo.Context, id string) (lesson Lesson, err error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return Lesson{}, err
	}

	query := datastore.NewQuery("lesson").Filter("ID =", id)

	var lessons Lessons

	keys, err := client.GetAll(ctx, query, &lessons)
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return Lesson{}, err
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>", keys)

	//対象ユーザのkeyを取得
	lessons[0].key = keys[0]

	//IDが重複することはないが、重複している場合は最新の講座を返す。
	return lessons[0], nil
}

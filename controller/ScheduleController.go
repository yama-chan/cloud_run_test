package controller

import (

	// "golang.org/x/net/context"
	// "github.com/labstack/echo/middleware"

	// "github.com/labstack/echo/v4"

	// "encoding/json"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/labstack/echo/v4"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
)

// userの構造体
type Schedule struct {
	key              *datastore.Key //データのkey ※更新時に使用
	ID               string         `json:"id"`
	User_ID          string         `json:"user_id"`
	StartDate        time.Time      `json:"start"`
	EndDate          time.Time      `json:"end"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Publishable      bool           `json:"publishable"`
	AllDay           bool           `json:"allDay"`
	Color            string         `json:"color"`
	RegisteredDate   string         `json:"registerddate"`
	LastModifiedDate string         `json:"lastmodifieddate"`
	Start_view       string         `json:"startdatetime"`
	End_view         string         `json:"enddatetime"`
}

type ViewScheduleStruct struct {
	ScheduleList []Schedule
	EditSchedule Schedule
	EditMode     bool
}

//userのスライスを型宣言
type Schedules []Schedule

var scheduleInf *Schedule

func GetScheduleJson(c echo.Context) []Schedule {
	fmt.Println("GetScheduleJson >>>>>>>>> called！！")

	var Schedules Schedules

	// event1 := Schedule{
	// 	StartDate: time.Date(2019, 10, 23, 1, 10, 6, 0, time.Local),
	// 	EndDate:   time.Date(2019, 10, 24, 1, 10, 6, 0, time.Local),
	// 	Title:     "topgate1",
	// }

	// Schedules = append(Schedules, event1)

	Schedules = GetSchedule(c)

	fmt.Println("schedules >>>>>>>>> ", Schedules)

	// // 配列をjsonに変換する
	// bytes, _ := json.Marshal(Schedules)
	// JsonStr := string(bytes)

	// fmt.Printf("json >>>>>>>", "%s\n", JsonStr)

	return Schedules
}

func GetScheduleViewStruct(c echo.Context) ViewScheduleStruct {

	//日付フォーマット
	const format = "2006/01/02 15:04:05" // 24h表現、0埋めあり

	viewScheduleStruct := ViewScheduleStruct{
		EditSchedule: Schedule{Color: "#3366CC"},
		ScheduleList: GetSchedule(c),
		EditMode:     false,
		// ScheduleJson: GetGetScheduleJson(c),
	}
	return viewScheduleStruct
}

func AddSchedule(c echo.Context) error {

	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	projectID := "sandbox-taisukeyamashita"

	log.Printf("ctx >>>>>>>>>>>>>>>>>>>>>>> %v", ctx)

	r := c.Request()
	r.ParseForm() //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。

	//日付フォーマット
	const format = "2006/01/02 15:04" // 24h表現、0埋めあり
	const format2 = "2006/01/02"      // 24h表現、0埋めあり

	const defaultColor = "#36c" // スケジュールのデフォルトの色分け

	//IDフォーマット
	const formatForID = "20060102150405"

	//モーダル画面にてスケジュール編集を行う
	// scheduleId := strings.Join(r.Form["mordal-scheduleId"], "")
	scheduleId := strings.Join(r.Form["scheduleId"], "")
	log.Printf("scheduleId >>>>>>>>>>>>>>>>>>>>>>> %v\n", scheduleId)

	//更新の場合
	if scheduleId != "" {
		log.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 更新\n")
		schedule, err := GetScheduleByID(c, scheduleId)
		scheduleInf = &schedule

		if err != nil {
			c.Render(http.StatusOK, "error", err.Error())
			return err
		}

		//新規登録の場合
	} else {
		log.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 新規\n")
		//保存対象となる構造体を定義,　ポイント型（&user）としないこと
		//《注意》：PUTの引数は構造体のポインタ(&user)とすること
		schedule := Schedule{}
		scheduleInf = &schedule

	}

	startDatetime := strings.Join(r.Form["startdate"], "")
	endDatetime := strings.Join(r.Form["enddate"], "")
	publishableFlg := strings.Join(r.Form["publishable"], "") == "on"

	// var start = time.Time{}
	// var end = time.Time{}

	start, err := time.Parse("2006/01/02 15:04", startDatetime)
	if err != nil {
		return err
	}

	end, err2 := time.Parse("2006/01/02 15:04", endDatetime)
	if err2 != nil {
		return err
	}

	fmt.Printf("startDatetime >>>>>>>> %v\n", startDatetime)
	fmt.Printf("endDatetime >>>>>>>> %v\n", endDatetime)
	fmt.Printf("publishableFlg >>>>>>>> %v\n", publishableFlg)
	fmt.Printf("start >>>>>>>> %v\n", start.Local())
	fmt.Printf("end >>>>>>>> %v\n", end.Local())
	fmt.Printf("err >>>>>>>> %v\n", err)
	fmt.Printf("err2 >>>>>>>> %v\n", err2)

	//保存対象となる構造体を定義,　ポイント型（&user）としないこと
	//《注意》：PUTの引数は構造体のポインタ(&schedule)とすること

	// //一意のID生成時処理のため排他ロック開始
	// schedule.mu.Lock()

	// //関数がreturnするまではロック
	// defer schedule.mu.Unlock()

	scheduleInf.Title = strings.Join(r.Form["title"], "")
	scheduleInf.StartDate = start

	//当日限りの予定はカレンダーに開始時刻のみ表示するため登録しない
	//ゆえにスケジュール詳細画面およびモーダルでは　End_view　を表示する。
	if start.Local().Format(format2) != end.Local().Format(format2) {
		scheduleInf.EndDate = end
	} else {
		//当日限りの予定の場合はEndDateにゼロ値を代入
		scheduleInf.EndDate = time.Time{}
	}
	scheduleInf.Publishable = publishableFlg
	scheduleInf.Description = strings.Join(r.Form["comment"], "")
	scheduleInf.Color = strings.Join(r.Form["color"], "")
	scheduleInf.Start_view = start.Local().Format(format)
	scheduleInf.End_view = end.Local().Format(format)
	// scheduleInf.RegisteredDate = registeredDate.Format(format)

	registeredDate := time.Now()

	//IDの無いユーザはIDを割り振って登録日時を更新する。
	if scheduleInf.ID == "" {
		//ID用の文字列スライスを生成
		//TODO: ID生成の方針として「ID　+ user_id」の組合せで一意のIDを生成する。
		s := []string{"SC-", registeredDate.Format(formatForID)}

		scheduleInf.RegisteredDate = registeredDate.UTC().Format(format)
		scheduleInf.ID = strings.Join(s, "")
		//すでにIDが割り振られているユーザは最終更新日時を変更する。
	} else {
		scheduleInf.LastModifiedDate = registeredDate.Format(format)
	}

	fmt.Println(r.Form) //これらのデータはサーバのプリント情報に出力されます
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return err
	}

	// Sets the kind for the new entity.
	kind := "schedule"
	// // Sets the name/ID for the new entity.
	// name := ""

	// keyインスタンスの生成
	// taskKey := datastore.NameKey(kind, name, nil)
	taskKey := datastore.IncompleteKey(kind, nil)

	fmt.Printf("taskKey >>>>>>>> %v\n", taskKey)
	fmt.Printf("client >>>>>>>> %v\n", client)

	//保存対象の構造体のポインタを引数に定義しPUT
	if scheduleInf.key != nil { //更新の場合　※エンティティがすでにkeyを持っている
		if _, err := client.Put(ctx, scheduleInf.key, scheduleInf); err != nil {
			s := []string{"Failed to save task:", err.Error()}
			c.Render(http.StatusOK, "error", strings.Join(s, ""))
			return err
		}

	} else { //新規の場合　※新規のkey(taskKey)を生成し保村処理を行う必要がある
		if _, err := client.Put(ctx, taskKey, scheduleInf); err != nil {
			s := []string{"Failed to save task:", err.Error()}
			c.Render(http.StatusOK, "error", strings.Join(s, ""))
			return err
		}
	}

	fmt.Printf("Saved %v: %v\n", taskKey, scheduleInf.Title)

	return nil
}

func GetScheduleByID(c echo.Context, id string) (schedule Schedule, err error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return Schedule{}, err
	}

	query := datastore.NewQuery("schedule").Filter("ID =", id)

	var schedules Schedules

	keys, err := client.GetAll(ctx, query, &schedules)
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return Schedule{}, err
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>", keys)

	//対象のレコードが存在しない場合
	if len(keys) == 0 {
		//独自エラーの作成
		err := errors.New("対象のレコードは存在しません。")
		c.Render(http.StatusOK, "error", err.Error())
		return Schedule{}, err
	}

	//対象ユーザのkeyをセット
	schedules[0].key = keys[0]

	//IDが重複することはないが、重複している場合は最新の講座を返す。
	return schedules[0], nil
}

func GetSchedule(c echo.Context) []Schedule {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		// TODO: Handle error.
	}

	var schedules Schedules

	// query := datastore.NewQuery("lesson").Filter("IsDeleted =", false)
	// query := datastore.NewQuery("schedule").Filter("Publishable =", true)
	query := datastore.NewQuery("schedule").Order("-Start_view")

	keys, err := client.GetAll(ctx, query, &schedules)
	if err != nil {
		// TODO: Handle error.
	}
	for i, key := range keys {
		fmt.Println(key)
		fmt.Println(schedules[i])
	}

	return schedules
}

func HandlePostAction(c echo.Context) (err error) {

	r := c.Request()
	r.ParseForm() //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。

	action := strings.Join(r.Form["action"], "")
	// action := r.Form["action"][0]
	if action == "add" {
		fmt.Println("form action >>>>>>>>>>>>>>>>>>>>>>>> ", "ADD")
		return AddSchedule(c)

	} else if action == "delete" {
		fmt.Println("form action >>>>>>>>>>>>>>>>>>>>>>>> ", "DELETE")
		return DeleteScheduleById(c)
	}
	newErr := errors.New("予期せぬリクエストを受信しました。")
	return newErr

}

func DeleteScheduleById(c echo.Context) (err error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return err
	}

	r := c.Request()
	r.ParseForm() //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。

	//スケジュールIDをhiddenから取得する。
	scheduleId := r.FormValue("scheduleId")
	log.Printf("scheduleId >>>>>>>>>>>>>>>>>>>>>>> %v\n", scheduleId)

	query := datastore.NewQuery("schedule").Filter("ID =", scheduleId)

	var schedules Schedules

	keys, err := client.GetAll(ctx, query, &schedules)
	if err != nil {
		// c.Render(http.StatusOK, "error", err.Error())
		return err
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>", keys)

	//対象のレコードが存在しない場合
	if len(keys) == 0 {
		//独自エラーの作成
		newErr := errors.New("削除対象のレコードは存在しませんでした。")
		// c.Render(http.StatusOK, "error", newErr.Error())
		return newErr
	}

	// 対象レコードの削除
	err = client.Delete(ctx, keys[0])

	if err != nil {
		//独自エラーの作成
		newErr := errors.New("対象レコードの削除にてエラーが発生しました。")
		// c.Render(http.StatusOK, "error", newErr.Error())
		return newErr
	}

	return nil
}

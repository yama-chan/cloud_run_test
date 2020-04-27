package controller

import (

	// "golang.org/x/net/context"
	// "github.com/labstack/echo/middleware"

	// "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
)

//一覧表示用の構造体
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
//※controllerパッケージ内ではpublicとなっている変数および型を自身のクラスで定義することなく使用できる。
//また、publicとなる変数名および型名は全て大文字始まりとなっていること。
type ViewIndexStruct struct {
	Users        []User
	Lessons      []Lesson
	ScheduleJson string
}

type LoginForm struct {
	UserId       string
	Password     string
	ErrorMessage string
}

// // lessonの構造体
// type Lesson struct {
// 	LessonName     string
// 	Communication  string
// 	Analytics      string
// 	Security       string
// 	Description    string
// 	RegisteredDate string
// }

// // //lessonのスライスを型宣言
// // type lessons []lesson

// // userの構造体
// type User struct {
// 	Fullname       string
// 	Sex            string
// 	Email          string
// 	Department     string
// 	RegisteredDate string
// 	IsDeleted      bool
// }

// //userのスライスを型宣言
// type users []user

func GetIndexViewStruct(c echo.Context) ViewIndexStruct {

	viewIndexStruct := ViewIndexStruct{
		Users:   GetUser2(c),
		Lessons: GetLesson(c),
		// ScheduleJson: GetGetScheduleJson(c),
	}
	return viewIndexStruct
}

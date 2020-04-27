package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	// "golang.org/x/net/context"
	// "github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"google.golang.org/appengine"

	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
	"context"

	controller "github.com/taisukeyamashita/test/controller"

	// controller "github.com/taisukeyamashita/test/controller"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	temp *template.Template
}

type templates *template.Template

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.temp.ExecuteTemplate(w, name, data)
}

type LoginForm struct {
	UserId       string
	Password     string
	ErrorMessage string
}

//lessonの構造体
type lesson struct {
	LessonName     string
	Communication  string
	Analytics      string
	Security       string
	Description    string
	RegisteredDate string
}

//lessonのスライスを型宣言
type lessons []lesson

//userの構造体
type user struct {
	Fullname       string
	Sex            string
	Email          string
	Department     string
	RegisteredDate string
	IsDeleted      bool
}

//userのスライスを型宣言
type users []user

var (
	// テンプレートディレクトリ
	templatesDir string = "view"
)

// Cookie型のstore情報
var store *sessions.CookieStore

func main() {

	//************ Echo ************ satrt
	// Echoのインスタンスを生成
	e := echo.New()

	//セッションを設定
	store := sessions.NewCookieStore([]byte("secret"))

	//セッション保持時間
	store.MaxAge(86400)

	//sessionオプション設定
	store.Options = &sessions.Options{
		Path: "/",
		// MaxAge=0 means no Max-Age attribute specified and the cookie will be
		// deleted after the browser session ends.
		// MaxAge<0 means delete cookie immediately.
		// MaxAge>0 means Max-Age attribute present and given in seconds.

		//ログイン(/login)画面遷移時にセッションを削除
		// MaxAge: -1,
		//７日間セッション(ログイン状態)を維持
		// MaxAge:   86400 * 7,
		//５分間セッション(ログイン状態)を維持
		MaxAge: 3000 * 7,

		HttpOnly: true,
	}

	//session
	e.Use(session.Middleware(store))

	renderer := &TemplateRenderer{
		temp: template.Must(template.ParseGlob("view/*.html")),
	}

	// テンプレートを利用するためのRendererの設定
	e.Renderer = renderer

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 静的ファイルのパスを設定
	e.Use(middleware.Static(""))
	e.Use(middleware.Static("/static"))

	e.GET("/", helloGET)
	e.GET("/login", helloGET)
	e.GET("/index", indexGET)
	e.GET("/editUser", editUserGET)
	// e.GET("/editUser/:Email", editUserGET)
	e.GET("/editLesson", editLessonGET)
	e.GET("/editSchedule", editScheduleGET)
	e.GET("/general", generalGET)
	e.GET("/table", tableGET)
	e.GET("/formComponent", formComponentGET)
	e.GET("/formValidation", formValidateGET)
	e.GET("/formButtons", formButtonsGET)
	e.GET("/formGrids", formGridsGET)
	e.GET("/widgets", widgetsGET)
	e.GET("/charts", chartsGET)
	e.GET("/profile", profileGET)
	e.GET("/contact", contactGET)
	e.GET("/blank", blankGET)
	e.GET("/getSchedule", getScheduleGET)
	e.GET("/addSchedule", editScheduleGET)
	e.GET("/addLesson", indexGET)
	e.GET("/addUser", indexGET)

	e.POST("/login", loginPOST)
	e.POST("/addUser", addUserPOST)
	e.POST("/addLesson", addLessonPOST)
	e.POST("/addSchedule", addSchedulePOST)
	e.POST("/editSchedule", editSchedulePOST)

	// サーバーをポート8080で起動
	http.Handle("/", e)
	appengine.Main()

	// e.Logger.Fatal(e.Start(":8080"))

	//************ Echo ************ end

	//************ golang ************ satrt

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// /
	http.HandleFunc("/", handleLogin)

	// /login
	http.HandleFunc("/login", handleLogin)

	// /index
	http.HandleFunc("/index", handleIndex)

	// /formValidation
	http.HandleFunc("/formValidation", handleFormValidation)

	// /formComponent
	http.HandleFunc("/formComponent", handleFormComponent)

	// /addUser
	http.HandleFunc("/editUser", handleEditUser)

	// /addUser
	http.HandleFunc("/addUser", handleAddUser)

	// サーバーをポート8080で起動
	//command:dev_appserver.py app.yaml

	//appengineで立ち上げる
	//cloud datastoreでinsertできるのはこっち
	// appengine.Main()

	//goの関数でwebサーバを立ち上げる
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// 	log.Printf("Defaulting to port %s", port)
	// }
	// log.Printf("Listening on port %s", port)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

	//************ golang ************ end

}

// 初期化を行います。
//echoインスタンスが初期化される前、main関数が実行される後に実行される
func init() {
	loadTemplates()
}

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存します（初期化時に実行）。
func loadTemplates() {
	fmt.Print("init() >>>>>>>> loadTemplates() called !! \n")
	// var baseTemplate = "view/index.html"
	// templates = make(map[string]*template.Template)
	// templates["index"] = template.Must(
	// 	template.ParseFiles(baseTemplate))
}

//test
func helloGET(c echo.Context) error {

	sess, _ := session.Get("session", c)
	// sess.Options = &sessions.Options{
	// 	Path: "/",
	// 	// MaxAge=0 means no Max-Age attribute specified and the cookie will be
	// 	// deleted after the browser session ends.
	// 	// MaxAge<0 means delete cookie immediately.
	// 	// MaxAge>0 means Max-Age attribute present and given in seconds.

	// 	//ログイン(/login)画面遷移時にセッションを削除
	// 	MaxAge: -1,

	// 	//７日間セッション(ログイン状態)を維持
	// 	// MaxAge:   86400 * 7,
	// 	HttpOnly: true,
	// }
	sess.Values["foo"] = "bar"
	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "login", LoginForm{})
}

func indexGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	fmt.Printf("session >>>>>>>>>>>>>>>>>>>>>>>> %v\n", sess.Values["loginCompleted"])

	return c.Render(http.StatusOK, "index", controller.GetIndexViewStruct(c))
}

func editUserGET(c echo.Context) error {

	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	// email := c.Param("Email") //プレースホルダEmailの値取り出し
	//URLパラメータをMap形式（map[string]　[]string）で取得する。
	urlValues := c.Request().URL.Query()

	//パラメータ名がidで複数指定されているURLに関しては、stringのスライスで格納される
	Id := urlValues["id"]
	fmt.Printf("param ID >>>>>>>>>>>>>>>>>>>>>>>> %v\n", Id)

	if len(Id) > 0 && Id[0] != "" {
		fmt.Printf("Id[0] >>>>>>>>>>>>>>>>>>>>>>>> %v\n", Id[0])
		//URLの最初に定義された値を検索条件とする。
		user, err := controller.GetUserByID(c, Id[0])
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "register_user", user)
	}
	return c.Render(http.StatusOK, "register_user", controller.User{})
}
func editLessonGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	//URLパラメータをMap形式（map[string]　[]string）で取得する。
	urlValues := c.Request().URL.Query()

	//パラメータ名がidで複数指定されているURLに関しては、stringのスライスで格納される
	Id := urlValues["id"]
	fmt.Printf("param ID >>>>>>>>>>>>>>>>>>>>>>>> %v\n", Id)

	if len(Id) > 0 && Id[0] != "" {
		fmt.Printf("Id[0] >>>>>>>>>>>>>>>>>>>>>>>> %v\n", Id[0])
		//URLの最初に定義された値を検索条件とする。
		user, err := controller.GetLessonByID(c, Id[0])
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "register_lesson", user)
	}
	return c.Render(http.StatusOK, "register_lesson", controller.Lesson{})
}

func editScheduleGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	//URLパラメータをMap形式（map[string]　[]string）で取得する。
	urlValues := c.Request().URL.Query()

	//パラメータ名がidで複数指定されているURLに関しては、stringのスライスで格納される
	Id := urlValues["id"]
	fmt.Printf("param ID >>>>>>>>>>>>>>>>>>>>>>>> %v\n", Id)

	if len(Id) > 0 && Id[0] != "" {
		fmt.Printf("Id[0] >>>>>>>>>>>>>>>>>>>>>>>> %v\n", Id[0])
		//URLの最初に定義された値を検索条件とする。
		schedule, err := controller.GetScheduleByID(c, Id[0])
		if err != nil {
			return err
		}

		ScheduleViewStruct := controller.GetScheduleViewStruct(c)
		ScheduleViewStruct.EditSchedule = schedule
		ScheduleViewStruct.EditMode = true

		return c.Render(http.StatusOK, "register_schedule", ScheduleViewStruct)
	}
	return c.Render(http.StatusOK, "register_schedule", controller.GetScheduleViewStruct(c))
}

func generalGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "general", "aaa")
}

func tableGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "basic_table", "aaa")
}
func formValidateGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "form_validation", "aaa")
}
func formComponentGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "form_component", "aaa")
}

func formButtonsGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "buttons", "aaa")
}
func formGridsGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "grids", "aaa")
}
func widgetsGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "widgets", "aaa")
}
func chartsGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "charts", "aaa")
}
func profileGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "profile", "aaa")
}
func contactGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "contact", "aaa")
}
func blankGET(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	return c.Render(http.StatusOK, "blank", "aaa")
}

func getScheduleGET(c echo.Context) error {

	fmt.Println("getScheduleGET >>>> 生成開始")
	return c.JSON(http.StatusOK, controller.GetScheduleJson(c))
}

func addUserPOST(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	//プロフィール画像選択時
	if c.Request().FormValue("action") == "uploadImage" {

		fmt.Println("action >>>>>>>>>>>>>>>>>>>>>")
		fmt.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 画像保存\n")

		const bucketName = "free_test_bucket"
		const filepath = "user/"

		user, err := controller.SaveImage(c, bucketName, filepath)
		if err != nil {
			return c.Render(http.StatusOK, "error", err.Error())
		}
		return c.Render(http.StatusOK, "register_user", user)
	}

	//登録処理
	err := controller.AddUser(c)

	if err != nil {

		return c.Render(http.StatusOK, "error", err.Error())
	}

	return c.Render(http.StatusOK, "index", controller.GetIndexViewStruct(c))
}
func addLessonPOST(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	//登録処理
	err := controller.AddLesson(c)

	if err != nil {

		return c.Render(http.StatusOK, "error", err.Error())
	}

	return c.Render(http.StatusOK, "index", controller.GetIndexViewStruct(c))
}

func addSchedulePOST(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	//登録処理
	err := controller.AddSchedule(c)

	if err != nil {

		return c.Render(http.StatusOK, "error", err.Error())
	}

	return c.Render(http.StatusOK, "register_schedule", controller.GetScheduleViewStruct(c))
}

func editSchedulePOST(c echo.Context) error {
	//セッション確認
	sess, sessErr := getSession(c)
	if sess == nil {
		return sessErr
	}

	//登録処理
	err := controller.HandlePostAction(c)

	if err != nil {
		return c.Render(http.StatusOK, "error", err.Error())
	}

	return c.Render(http.StatusOK, "register_schedule", controller.GetScheduleViewStruct(c))
}

func loginPOST(c echo.Context) error {

	//ruquestを取得
	r := c.Request()
	r.ParseForm()

	userId := strings.Join(r.Form["UserId"], "")
	password := strings.Join(r.Form["Password"], "")

	//登録処理
	loginForm := LoginForm{
		UserId:   userId,
		Password: password,
	}

	fmt.Printf("userId >>>> %v\n", loginForm.UserId)
	fmt.Printf("password >>>> %v\n", loginForm.Password)

	// userId := html.EscapeString(loginForm.UserId)
	// password := html.EscapeString(loginForm.Password)

	if userId != "userId" && password != "password" {
		loginForm.ErrorMessage = "ユーザーID または パスワードが間違っています"
		return c.Render(http.StatusOK, "login", loginForm)
	}

	//セッションにデータを保存する
	// session := session.Default(c)
	sess, _ := session.Get("session", c)
	sess.Values["loginCompleted"] = "completed"
	sess.Save(c.Request(), c.Response())

	fmt.Printf("controller.GetIndexViewStruct(c) >>>>>>>>>>>>>>>> %v\n", controller.GetIndexViewStruct(c))
	// fmt.Printf("controller.GetIndexViewStruct(c) users >>>>>>>>>>>>>>>> %v\n", controller.GetIndexViewStruct(c).users)
	// fmt.Printf("controller.GetIndexViewStruct(c) lessons >>>>>>>>>>>>>>>> %v\n", controller.GetIndexViewStruct(c).lessons)

	return c.Render(http.StatusOK, "index", controller.GetIndexViewStruct(c))

	// return c.Render(http.StatusOK, "index", controller.GetUser2(c))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// テンプレートをパース
	tpl := template.Must(template.ParseFiles("view/login.html"))

	// テンプレートを描画
	tpl.Execute(w, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

	//user構造体の取得条件(IsDeleted = false)のユーザのみ
	user := user{
		IsDeleted: false}

	ctx := appengine.NewContext(r)

	fmt.Printf("users >>>> %v\n", user.getUser(ctx))

	execTemplate("index", w, user.getUser(ctx), "index", "sidebar", "header", "credits")

}

func handleAddUser(w http.ResponseWriter, r *http.Request) {

	// GET or POST で処理の分岐

	//GET
	if r.Method == http.MethodGet {

		// テンプレートをパース
		tpl := template.Must(template.ParseFiles("view/register_user.html"))
		// addUser(w, r)

		// テンプレートを描画
		tpl.Execute(w, nil)

		//POST
	} else if r.Method == http.MethodPost {

		if !controller.AddRegisteredUser(r) {
			tpl := template.Must(template.ParseFiles("view/error.html"))
			// テンプレートを描画
			tpl.Execute(w, nil)

			return
		}

		io.WriteString(w, "This is a post request completed!")

	} else {
		io.WriteString(w, "This is undefined")

	}
}

func handleFormValidation(w http.ResponseWriter, r *http.Request) {
	// テンプレートをパース
	tpl := template.Must(template.ParseFiles("view/form_validation.html"))

	// テンプレートを描画
	tpl.Execute(w, nil)
}

func handleFormComponent(w http.ResponseWriter, r *http.Request) {
	// テンプレートをパース
	tpl := template.Must(template.ParseFiles("view/form_component.html"))

	// テンプレートを描画
	tpl.Execute(w, nil)
}

func handleEditUser(w http.ResponseWriter, r *http.Request) {

	// GET or POST で処理の分岐

	//GET
	if r.Method == http.MethodGet {

		execTemplate("register_user", w, nil, "register_user", "sidebar", "header", "credits")

		//POST
	} else {
		io.WriteString(w, "This is undefined")

	}
}

// execTemplate はfilesのテンプレートからHTMLを構築して,
// wに対して書き込みます.
// HTML構築の際にはdataを利用します.
// filesで指定するテンプレートには必ず{{define "layout"}}された
// ものを1つだけ含む必要があります.
func execTemplate(layoutFile string, w http.ResponseWriter, data interface{}, files ...string) {
	// 渡された引数からテンプレートパスの集合を作る
	var pathes []string
	for _, f := range files {
		p := fmt.Sprintf("%s/%s.html", templatesDir, f)
		fmt.Printf("p >>>>>>>>>>>>>>>>>>>> %v\n", p)
		pathes = append(pathes, p)
	}
	//上記で作ったパスの一覧を使ってテンプレートを作る
	template := template.Must(template.ParseFiles(pathes...))
	// layoutが必ず基点になるという事にする
	template.ExecuteTemplate(w, layoutFile, data)
}

func getTemplate(files ...string) *TemplateRenderer {
	// 渡された引数からテンプレートパスの集合を作る
	var pathes []string
	for _, f := range files {
		p := fmt.Sprintf("%s/%s.html", templatesDir, f)
		fmt.Printf("p >>>>>>>>>>>>>>>>>>>> %v\n", p)
		pathes = append(pathes, p)
	}
	//上記で作ったパスの一覧を使ってテンプレートを作る
	renderer := &TemplateRenderer{
		temp: template.Must(template.ParseFiles(pathes...)),
	}

	return renderer
}

func (user *user) getUser(ctx context.Context) []user {

	// クエリー作成
	// q := datastore.NewQuery("user").Filter("IsDeleted =", user.IsDeleted).Filter("Fullname =", p.Fullname)
	// q := datastore.NewQuery("user").Filter("IsDeleted =", false)

	//ローカル変数を定義（上記でtype宣言した型で宣言）
	var users users

	// //クエリーからGerAllでまとめてレコードを取得
	// _, err := q.GetAll(ctx, &users)
	// if err != nil {
	// 	log.Fatalf("Error fetching next task: %v", err)
	// }

	return users

}

//セッションを取得
func getSession(c echo.Context) (*sessions.Session, error) {

	// //独自エラーの作成
	// newErr := errors.New("セッションエラー：")

	sess, err := session.Get("session", c)
	if err != nil {
		return nil, c.Render(http.StatusOK, "error", err.Error())
	}
	if sess.Values["loginCompleted"] != "completed" {
		return nil, c.Render(http.StatusOK, "login", LoginForm{})
	}
	return sess, nil
}

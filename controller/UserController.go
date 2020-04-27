package controller

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
	"cloud.google.com/go/storage"

	//ユニークID生成用
	"github.com/rs/xid"
)

//Userの構造体
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type User struct {
	mu               sync.Mutex     //排他制御用　※一意のIDを生成するためユーザ同時登録を防ぐ
	key              *datastore.Key //データのkey ※更新時に使用
	ID               string
	Fullname         string
	Sex              string
	Email            string
	Department       string
	Description      string
	RegisteredDate   string
	LastModifiedDate string
	IsDeleted        bool
	IconURL          string
}

//userのスライスを型宣言
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type Users []User

var userInf *User

func AddUser(c echo.Context) error {

	// Set your Google Cloud Platform project ID.
	const projectID = "sandbox-taisukeyamashita"
	const bucketName = "free_test_bucket"
	const filepath = "user/"

	// if c.Request().FormValue("action") == "uploadImage" {
	// 	fmt.Println("action >>>>>>>>>>>>>>>>>>>>>")
	// 	fmt.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 画像保存\n")
	// 	user, err := SaveImage(c, bucketName, filepath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	c.Render(http.StatusOK, "register_user", user)
	// 	return err
	// }

	userInf := User{}

	userInf, err1 := SaveImage(c, bucketName, filepath)

	if err1 != nil {
		s := []string{"Failed to save Image:", err1.Error()}
		c.Render(http.StatusOK, "error", strings.Join(s, ""))
		return err1
	}

	err2 := saveUser(c, projectID, userInf.IconURL)

	// err = saveImage(r, ctx)

	if err2 != nil {
		s := []string{"Failed to save User:", err2.Error()}
		c.Render(http.StatusOK, "error", strings.Join(s, ""))
		return err2
	}

	fmt.Printf("Saved : %v\n", userInf.Fullname)

	return nil
}

func AddRegisteredUser(r *http.Request) bool {

	//golang 1.11使用時
	ctx := context.Background()

	//golang 1.12 使用時
	// ctx := appengine.NewContext(r)

	log.Printf("ctx >>>>>>>>>>>>>>>>>>>>>>> %v", ctx)

	r.ParseForm() //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。

	fmt.Println(r.Form) //これらのデータはサーバのプリント情報に出力されます
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	//r.Form["fullname"]ではvalueを配列[]で取得するため size=1 でもrangeまたは配列要素を結合するなどの処理が必要となる
	fmt.Println(r.Form["fullname"])
	fmt.Println(r.Form["email"])
	fmt.Println(r.Form["department"])
	fmt.Println(r.Form["sex"])

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		fmt.Println("val:", v)
	}

	//日付フォーマット
	const format = "2006/01/02 15:04:05" // 24h表現、0埋めあり

	//保存対象の構造体を定義
	user := User{
		Fullname:       r.FormValue("fullname"),
		Sex:            r.FormValue("sex"),
		Email:          r.FormValue("email"),
		Department:     r.FormValue("department"),
		Description:    r.FormValue("comment"),
		IsDeleted:      false,
		RegisteredDate: time.Now().Format(format)}

	log.Print(user)

	// //自動的にエンティティのキーとして数値IDを取得
	// key := datastore.NewIncompleteKey(ctx, "user", nil)
	// log.Printf("key >>>>>>>>>>>>>>>>>>>>>>> %v", key)

	// if _, err := datastore.Put(ctx, key, user); err != nil {

	// 	//エラー画面に遷移させる予定だが一旦やめる・・
	//  //保存対象の構造体のポインタを引数に定義しPUT
	// 	// if _, err := datastore.Put(nil, nil, &user); err != nil {
	// 	// tpl := template.Must(template.ParseFiles("../view/404.html"))
	// 	// テンプレートを描画
	// 	// tpl.Execute(w, nil)
	// 	log.Printf("err >>>>>>>>>>>>>>>>>>>>>>> %v", err)

	// 	return false

	// }

	return true
}

func GetUser2(c echo.Context) []User {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		// TODO: Handle error.
	}

	var users Users

	query := datastore.NewQuery("user").Filter("IsDeleted =", false)

	keys, err := client.GetAll(ctx, query, &users)
	if err != nil {
		// TODO: Handle error.
	}
	for i, key := range keys {
		fmt.Println(key)
		fmt.Println(users[i])
	}

	return users
}

func GetUserByID(c echo.Context, id string) (user User, err error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "sandbox-taisukeyamashita")
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return User{}, err
	}

	query := datastore.NewQuery("user").Filter("ID =", id)

	var users Users

	keys, err := client.GetAll(ctx, query, &users)
	if err != nil {
		c.Render(http.StatusOK, "error", err.Error())
		return User{}, err
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>", keys)

	//クエリーの結果が0（nil） の場合
	//IDは生成されているがDataStoreに登録されていない場合は空のuserを返す
	if keys == nil {
		return User{}, nil
	}

	//対象講座のkeyを取得
	users[0].key = keys[0]

	for i, key := range keys {
		fmt.Println("keys[]  >>>>>>", key, " : ", keys[i])
	}

	//IDが重複することはないが、重複している場合は最新のユーザを返す。
	return users[0], nil
}

//ファイルのcontent-typeを返す
func GetFileContentType(header *multipart.FileHeader) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	byteContainer := make([]byte, 512)

	// file is a *multipart.FileHeader gotten from http request.
	// ファイル本体はメモリまたはディスクに格納されているため、*FileHeader の Openメソッドにてファイルにアクセスする.
	fileContent, err := header.Open()
	if err != nil {
		return "", err
	}

	// Read content of *multipart.FileHeader into []byte
	// ファイルを読み込みbyteスライスに格納
	fileContent.Read(byteContainer)
	fmt.Println(byteContainer)

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.

	contentType := http.DetectContentType(byteContainer)
	fmt.Println("content-type >>>>>>>", contentType)

	return contentType, nil
}

// Cloud Storageに画像ファイルを保存する。aaaaaaaabb
// Formの入力内容およびCloudStrorageに保存した画像の公開URLを格納したuserの構造体を返す。
func SaveImage(c echo.Context, bucketName string, path string) (user User, err error) {

	ctx := context.Background()
	log.Printf("ctx >>>>>>>>>>>>>>>>>>>>>>> %v", ctx)

	//echo.Contextから*http.Requestを取得
	r := c.Request()

	//urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。
	r.ParseForm()

	//現時点でのフォームの入力内容をセット
	userForm := User{
		ID:          r.FormValue("userId"),
		Fullname:    r.FormValue("fullname"),
		Sex:         r.FormValue("sex"),
		Email:       r.FormValue("email"),
		Department:  r.FormValue("department"),
		Description: r.FormValue("comment"),
		// IconURL:     fmt.Sprintf(publicURL, bucketName, filename),
	}

	// 画像ファイルの取得
	file, fileHeader, err := r.FormFile("upload")

	if file == nil && fileHeader == nil {
		return User{}, nil
	}

	log.Println("file >>>>>>>>>>>>>>>>>>>>>>> ", file)
	log.Println("fileHeader >>>>>>>>>>>>>>>>>>>>>>> ", fileHeader)

	if err != nil {
		// s := []string{"ファイルの取得処理にてエラーが発生しました。:", err.Error()}
		// c.Render(http.StatusOK, "error", strings.Join(s, ""))
		return User{}, err
	}
	defer file.Close()

	// 取得した画像のcontent-typeを取得
	contentType, err := GetFileContentType(fileHeader)
	if err != nil {
		// s := []string{"content-typeを取得時にエラーが発生しました。:", err.Error()}
		// c.Render(http.StatusOK, "error", strings.Join(s, ""))
		return User{}, err
	}

	log.Println("contentType >>>>>>>>>>>>>>>>>>>>>>> ", contentType)

	// 取得したcontentTypeから{image/}を空文字に書き換えて拡張子を取得する。
	extension := strings.Replace(contentType, "image/", "", -1)
	fmt.Printf("extension >>>>>> : %s\n", extension)

	// 画像ファイルの保存準備
	gsclient, err := storage.NewClient(ctx) // クライアント作成
	if err != nil {
		// s := []string{"Cloud Storageクライアント作成エラー:", err.Error()}
		// c.Render(http.StatusOK, "error", strings.Join(s, ""))
		return User{}, err
	}
	defer gsclient.Close()

	if userForm.ID == "" {
		userForm.ID = createUniqueId()
	}

	// 画像ファイルの保存先のフォルダ構成とファイル名の設定
	filenameStr := []string{path, userForm.ID, ".", extension}
	filename := strings.Join(filenameStr, "")
	log.Println("filename >>>>>>>>>>>>>>>>>>>>>>> ", filename)

	// 保存先のバケットとファイルを設定
	bucket := gsclient.Bucket(bucketName) // 接続先バケット
	outObj := bucket.Object(filename)     // アップロード先オブジェクト
	writer := outObj.NewWriter(ctx)       // アップロードするためのライター
	// 上記の処理をまとめて記述するならメソッドチェーンで繋げて以下のように記述するのが一般的？かな
	// writer := client.Bucket(bucket).Object(object).NewWriter(ctx)

	// ACLを「すべてのユーザー」が「読み取り」できるように設定
	writer.ACL = []storage.ACLRule{
		{Entity: storage.AllUsers, Role: storage.RoleReader}}
	// キャッシュの設定
	// writer.CacheControl = "public, max-age=86400" //一日
	writer.ObjectAttrs.CacheControl = "no-cache" //キャッシュなし

	//重要なポイントは writer.Close() にある。実は、サーバーにファイルをアップロードする部分はここである。
	defer writer.Close()

	// content-TypeをHeaderから取得する。上記のメソッド（GetFileContentType）で取得する方法より楽な方法かもしれない。。。
	// *** また、下のように ContentType を具体的に指定するのは良くない。 ***
	//  **  writer.ContentType = "image/png"  **
	writer.ContentType = fileHeader.Header.Get("Content-Type")

	// 画像ファイルをwriterにコピー
	if _, err = io.Copy(writer, file); err != nil {
		return User{}, err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"

	userForm.IconURL = fmt.Sprintf(publicURL, bucketName, filename)

	return userForm, nil
}

func saveUser(c echo.Context, projectID string, iconURL string) (err error) {

	ctx := context.Background()
	log.Printf("ctx >>>>>>>>>>>>>>>>>>>>>>> %v", ctx)

	//echo.Contextから*http.Requestを取得
	r := c.Request()

	//urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。
	r.ParseForm()

	//日付フォーマット
	const format = "2006/01/02 15:04:05" // 24h表現、0埋めあり

	userId := r.FormValue("userId")
	fmt.Printf("userId >>>>>>>>>>>>>>>>>>>>>>> %v\n", userId)

	//更新の場合
	if userId != "" {
		fmt.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 更新\n")
		user, err := GetUserByID(c, userId)
		if err != nil {
			// c.Render(http.StatusOK, "error", err.Error())
			return err
		}
		if user.ID != "" {
			// ユーザIDがすでに割り振られておりDatastoreにもデータがある場合
			userInf = &user
		} else {
			// ユーザIDがすでに割り振られているがDatastoreにまだ登録されていない場合はIDのみ登録
			userInf = &User{ID: userId}
		}

		//新規登録の場合
	} else {
		fmt.Printf("action >>>>>>>>>>>>>>>>>>>>>>> 新規\n")
		//保存対象となる構造体を定義,　ポイント型（&user）としないこと
		//《注意》：PUTの引数は構造体のポインタ(&user)とすること
		user := User{}
		userInf = &user
	}

	//一意のID生成時処理のため排他ロック開始
	userInf.mu.Lock()

	//関数がreturnするまではロック
	defer userInf.mu.Unlock()

	//登録日時秒の組み合わせは一意とするためロック内で定義
	registeredDate := time.Now()

	// userInf.Fullname = strings.Join(r.PostForm["fullname"], "")
	userInf.Fullname = r.FormValue("fullname")
	userInf.Sex = r.FormValue("sex")
	userInf.Email = r.FormValue("email")
	userInf.Department = r.FormValue("department")
	userInf.Description = r.FormValue("comment")
	userInf.IconURL = iconURL
	// userInf.IconURL = r.FormValue("iconURL")iconURL
	userInf.IsDeleted = false
	// userInf.RegisteredDate = registeredDate.Format(format)

	//IDの無いユーザは新規IDを割り振って登録日時を更新する。
	if userInf.ID == "" {
		//ID用の文字列スライスを生成
		s := []string{"U-", createUniqueId()}

		userInf.RegisteredDate = registeredDate.Format(format)
		userInf.ID = strings.Join(s, "")
		//すでにIDが割り振られているユーザは最終更新日時を変更する。
	} else {
		userInf.LastModifiedDate = registeredDate.Format(format)
	}

	// ログ
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	//r.Form["fullname"]ではvalueを配列[]で取得するため size=1 でもrangeまたは配列要素を結合するなどの処理が必要となる
	fmt.Println(r.FormValue("fullname"))
	fmt.Println(r.Form["email"])
	fmt.Println(r.Form["department"])
	fmt.Println(r.Form["sex"])
	fmt.Println(r.Form["iconURL"])

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		// s := []string{"Failed to create client:", err.Error()}
		// c.Render(http.StatusOK, "error", strings.Join(s, ""))
		return err
	}

	// 新規エンティティのkindを設定
	kind := "user"
	// Sets the name/ID for the new entity.
	// name := ""

	// keyインスタンスの生成
	// taskKey := datastore.NameKey(kind, name, nil)
	taskKey := datastore.IncompleteKey(kind, nil)

	fmt.Printf("taskKey >>>>>>>> %v\n", taskKey)
	fmt.Printf("client >>>>>>>> %v\n", client)

	//保存対象の構造体のポインタを引数に定義しPUT
	if userInf.key != nil { //更新の場合　※エンティティがすでにkeyを持っている場合
		if _, err := client.Put(ctx, userInf.key, userInf); err != nil {
			// s := []string{"Failed to save task:", err.Error()}
			// c.Render(http.StatusOK, "error", strings.Join(s, ""))
			return err
		}

	} else { //新規の場合　※新規のkey(taskKey)を生成し保村処理を行う必要がある場合
		if _, err := client.Put(ctx, taskKey, userInf); err != nil {
			// s := []string{"Failed to save task:", err.Error()}
			// c.Render(http.StatusOK, "error", strings.Join(s, ""))
			return err
		}
	}

	return nil
}

func createUniqueId() (id string) {
	// ************************
	// xidについて
	// ************************
	// 詳しくはGitHubのREADMEの書かれていますが、その中から一部抜粋して紹介します。

	// binaryのformat
	// 全体で12bytesで、先頭から以下のように構成されています。

	// 4bytes: Unix timestamp (秒単位)
	// 3bytes: ホストの識別子
	// 2bytes: プロセスID
	// 3bytes: ランダムな値からスタートしたカウンタの値

	// 生成される文字列
	// 20文字のlower caseの英数字。([0-9a-v]{20})
	// 例: b8hpcg8hv3amvi9dol0g

	// idを生成
	guid := xid.New()
	fmt.Println(guid.String())

	// binaryの各partの情報（参考に）
	machine := guid.Machine()
	pid := guid.Pid()
	time := guid.Time()
	counter := guid.Counter()
	fmt.Printf("machine: %v, pid: %v, time: %v, counter: %v\n", machine, pid, time, counter)

	return guid.String()
}

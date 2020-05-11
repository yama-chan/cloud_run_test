package config

import (
	"sync"
)

// アプリケーションの関心事に関する設定情報
var (
	//固定の設定情報
	CloudRunServiceName    string = "test_golang_application"
	GCS_PublicURL_Template string = "https://storage.googleapis.com/%s/%s" //バケット名/ファイルパス

	//動的に設定する設定値
	ProjectID         string
	CsvBucketFilePath string //"csv/"
	IndexFileName     string //"index.html"
	MyFreeBucketName  string //"free_test_bucket"
	TemplateBucket    string //templateの格納バケット名 ※環境変数で設定
	TemplateDir       string //templateの格納ディレクトリ ※環境変数で設定
)

//このファイルでは一度だけ設定情報の初期化処理を実行する
var initializeOnce sync.Once

// Initialize 設定情報を初期化
func Initialize(
	projectID string,
	csvBucketFilePath string,
	myFreeBucketName string,
	templateBucket string,
	templateDir string,
	indexFileName string,
) {
	initializeOnce.Do(func() {
		ProjectID = projectID
		CsvBucketFilePath = csvBucketFilePath
		MyFreeBucketName = myFreeBucketName
		TemplateBucket = templateBucket
		TemplateDir = templateDir
		IndexFileName = indexFileName
	})
}

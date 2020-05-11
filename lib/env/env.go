package env

import (
	"os"
	"sync"
	"time"

	"github.com/taisukeyamashita/test/lib/config"
)

// envCloudRun アプリケーションの環境設定を行う構造体
type envCloudRun struct {
	config         InitializeConfig
	initializeOnce sync.Once
}

var _ EnvValues = new(envCloudRun)

// InitializeConfig appengin用コンフィグ
type InitializeConfig struct {
	DATASTORE     DatastoreConfig
	GCS           GCSConfig
	Port          string
	TemplateDir   string //templateの格納ディレクトリ ※環境変数で設定
	IndexFileName string //"index.html"
}

// Cloud Runのデフォルト環境変数を取得
var (
	// projectID       = os.Getenv("GOOGLE_CLOUD_PROJECT")
	CloudRunServiceName     = os.Getenv("K_SERVICE")
	CloudRunServiceRevision = os.Getenv("K_REVISION")
	CloudRunServiceConfig   = os.Getenv("K_CONFIGURATION")
)

// DatastoreConfig DATASTOREの設定情報
type DatastoreConfig struct {
	ProjectID string
}

// GCSConfig GCSの設定情報
type GCSConfig struct {
	MyFreeBucketName  string //"free_test_bucket"
	CsvBucketFilePath string //"csv/"
	TemplateBucket    string //templateの格納バケット名 ※環境変数で設定
}

// CreateInitializeConfig 環境変数やデフォルト値からConfig値を設定及び取得する
func CreateInitializeConfig() InitializeConfig {
	return InitializeConfig{
		DATASTORE: DatastoreConfig{
			ProjectID: MustGetString("GOOGLE_CLOUD_PROJECT"),
		},
		GCS: GCSConfig{
			MyFreeBucketName:  GetEnvString("MY_FREE_BUCKET_NAME", ""),
			CsvBucketFilePath: GetEnvString("CSV_BUCKET_FILE_PATH", ""),
			TemplateBucket:    GetEnvString("TEMPLATE_BUCKET", ""),
		},
		Port:          MustGetString("PORT"),
		TemplateDir:   GetEnvString("TEMPLATE_DIR", ""),
		IndexFileName: GetEnvString("INDEX_FILE_NAME", ""),
	}
}

var (
	// アプリケーションの環境情報をポインタでファイル変数に格納
	applicationEnvValues *envCloudRun
	// このファイルのMutex
	applicationEnvValuesMux = new(sync.Mutex)
)

// GetEnvValues リソースを取得する
func GetEnvValues(config InitializeConfig) EnvValues {
	// ファイルをロック
	applicationEnvValuesMux.Lock()
	defer applicationEnvValuesMux.Unlock()

	if applicationEnvValues == nil {
		// ファイルに格納
		applicationEnvValues = &envCloudRun{config: config}
	}
	return applicationEnvValues
}

// Initialize リソースの初期化 implement EnvValues interface
func (env *envCloudRun) Initialize() {
	env.initializeOnce.Do(func() {
		// アプリケーションの関心事に関する設定情報(config.go)を初期化
		env.initializeConfig()
		// UTCになるので明示的にJST変換する
		time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	})
}

// initializeConfig アプリケーションの関心事に関する設定情報(config.go)を初期化
func (env *envCloudRun) initializeConfig() {
	config.Initialize(
		env.config.DATASTORE.ProjectID,
		env.config.GCS.CsvBucketFilePath,
		env.config.GCS.MyFreeBucketName,
		env.config.GCS.TemplateBucket,
		env.config.TemplateDir,
		env.config.IndexFileName,
	)
}

// Port ポート番号を取得する implement EnvValues interface
func (env *envCloudRun) Port() string {
	return env.config.Port
}

// Finalize リソースの開放処理 implement EnvValues interface
func (env *envCloudRun) Finalize() {
	applicationEnvValuesMux.Lock()
	defer applicationEnvValuesMux.Unlock()
	// 環境情報を開放
	applicationEnvValues = nil
}

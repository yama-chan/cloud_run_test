package env

// EnvValues アプリケーションの環境情報(周り)の設定を行うインタフェースを用意
type EnvValues interface {
	Initialize()
	Port() string
	Finalize()
}

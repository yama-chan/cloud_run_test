package env

import (
	"fmt"
	"os"
	"strconv"
)

// MustGetString 文字列の環境変数を取得する。
// 対象変数がない場合はpanicを発生させる
func MustGetString(key string) string {
	v := GetEnvString(key, "")
	if v == "" {
		panic(fmt.Sprintf("%s is blank", key))
	}
	return v
}

// GetEnvString 文字列の環境変数を取得する対象変数がない場合はdefaultValueを返す
func GetEnvString(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultValue
}

// GetEnvUnsignedInt 数値の環境変数を取得する取得する対象変数がない場合、defaultValueを返す。
func GetEnvUnsignedInt(key string, defaultValue uint32) uint32 {
	v := os.Getenv(key)
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return uint32(i)
	}
	return defaultValue
}

// GetEnvUnsignedInt64 64bit数値の環境変数を取得する取得する対象変数がない場合、defaultValueを返す。
func GetEnvUnsignedInt64(key string, defaultValue uint64) uint64 {
	v := os.Getenv(key)
	if v != "" {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		return i
	}
	return defaultValue
}

// GetEnvBool bool値の環境変数を取得する取得する対象変数がない場合、defaultValueを返す。
func GetEnvBool(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	value, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}
	return value
}

// OnLocalDevServer ローカルで実行されているかの判定値
var OnLocalDevServer = GetEnvBool("DEV_SERVER", false)

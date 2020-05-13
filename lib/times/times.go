package times

import (
	"time"
)

// CurrentTime 現在時刻を返す
//
// テストなどで現在時刻などを任意に変更可能なように、グローバル変数として定義
var CurrentTime = func() time.Time {
	return time.Now()
}

// ローカル時間を初期化
func SetJSTTime() {
	// UTCになるので明示的にJST変換する
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
}

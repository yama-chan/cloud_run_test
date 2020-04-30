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

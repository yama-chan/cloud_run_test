package mail

import (
	"context"
	"time"
)

// EmailContentType EmailContentType
type (
	EmailContentType string
	// TODO: 型は別途作成するようにする。typeパッケージ作る？
	// EmailAddress     string
)

// Email区分
const (
	HTMLEmail EmailContentType = "text/html"
	TextEmail EmailContentType = "text/plain"
)

type (
	// EmailSender メール送信インターフェース
	EmailSender interface {
		SendBroad(ctx context.Context, conf EmailSenderConfig) error
	}

	// EmailSenderConfig メール送信設定
	EmailSenderConfig struct {
		Message   Message
		FromName  string
		FromEmail string
		// ※トランザクショナルテンプレートを使用してリクエストにtemplate_IDを指定している場合はContentは必須ではありません。
		// ContentType EmailContentType
		To []string
		// TODO: Goルーチンを使って複数メール送信の際に並列で送信出来るか実装する
		// SendFailedChan chan<- EmailSendFailed
	}

	// Message メールの内容
	Message struct {
		// ID string `validate:"required"` // お知らせID
		// Type         NotificationType `validate:"required"` // お知らせタイプ
		RegisteredAt time.Time `validate:"required"` // 作成日時
		Title        string    `validate:"required"` // タイトル
		// Body         string    // メール本文 ※テンプレートを使用するので今回はコメント
	}
)

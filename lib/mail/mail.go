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
		// To []EmailAddress
		To []string
		// SendFailedChan chan<- EmailSendFailed
	}

	Message struct {
		ID string `validate:"required"` // お知らせID
		// Type         NotificationType `validate:"required"` // お知らせタイプ
		Title        string    `validate:"required"` // タイトル
		Body         string    // 詳細内容
		RegisteredAt time.Time `validate:"required"` // 作成日時
	}
)

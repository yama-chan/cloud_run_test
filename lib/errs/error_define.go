package errs

import (
	"fmt"
	"net/http"
)

var (
	BadRequest   = newDefaultErr(2001, http.StatusBadRequest, "正しくないリクエストが送られました")
	Unauthorized = newDefaultErr(2002, http.StatusUnauthorized, "認証に失敗しました")
	Forbidden    = newDefaultErr(2003, http.StatusForbidden, "アクセス権限がありません")
	NotFound     = newDefaultErr(2004, http.StatusNotFound, "対象データを見つけられませんでした")
	Conflict     = newDefaultErr(2005, http.StatusConflict, "データが競合しています")
	ServerError  = newDefaultErr(3001, http.StatusInternalServerError, "不明なエラーが発生しました")
)

// AppError errorインターフェースを実装したインターフェース
// interface(今回はerrorインターフェース)を埋め込むことでインタフェースを満たす
type AppError interface {
	error //error interface { Error() string }なので、AppErrorをimplementする場合はError() stringを実装する
	StatusCode() int
	WithMessage(message string) AppError
}

// DefaultErr implement AppError
type DefaultErr struct {
	msg        string
	code       int
	statusCode int
}

var _ AppError = &DefaultErr{}

func newDefaultErr(code int, statusCode int, msg string) *DefaultErr {
	return &DefaultErr{
		msg:        msg,
		code:       code,
		statusCode: statusCode,
	}
}

// DefaultErr implement error
func (w *DefaultErr) Error() string {
	return fmt.Sprintf("[error] %s", w.msg)
}

// DefaultErr implement AppError
func (w *DefaultErr) StatusCode() int {
	return w.statusCode
}

// DefaultErr implement AppError
func (w *DefaultErr) WithMessage(message string) AppError {
	// errMsg := xerrors.Errorf("%s\n: %v",message, w).Error()
	w.msg = message
	return w
}

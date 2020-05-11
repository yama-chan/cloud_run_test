package errs

import (
	"fmt"

	"golang.org/x/xerrors"
)

// WrapXerror 既存のエラーから新規のエラーを作成する
func WrapXerror(err error) error {
	if err == nil {
		return nil
	}
	xerr := xerrors.Errorf("[wrapped error]: %w", err)
	return xerr
}

func NewXerrorWithMessage(message string) error {
	return xerrors.New(fmt.Sprintf("[error]: %s", message))
}

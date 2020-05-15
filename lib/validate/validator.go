package validate

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/taisukeyamashita/test/lib/errs"
)

// TODO: バリデート処理の充実
// https://godoc.org/gopkg.in/go-playground/validator.v10
var (
	// ファイルが読み込まれるたび（バリデーションが実行されるたびに）に、変数'vd' *validator.Validate を初期化する
	vd = func() *validator.Validate {
		v := validator.New()
		v.RegisterValidation("date", dateValidation)
		v.RegisterValidation("datetime", datetimeValidation)
		// v.RegisterValidation("date_range", dateRangeValidation) // TODO: モデル独自で　Validate関数を持ってもいいかも
		log.Printf("validator is initialized")
		return v
	}()
)

// 構造体のフィールドバリデート
// ValidateStruct validates a structs exposed fields, and automatically validates nested structs
func ValidateStruct(v interface{}) error {
	// TODO: フィールドにtime.Time型がある場合のバリデートはどうなるのか? 調査・検証する
	if err := vd.Struct(v); err != nil {
		return errs.NewXerrorWithMessage(fmt.Sprintf("[field Error]: %v", err.Error()))
	}
	return nil
}

func dateValidation(fl validator.FieldLevel) bool {

	_, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return true
}

func datetimeValidation(fl validator.FieldLevel) bool {
	// TODO: datetimeのフォーマットを決める、できればアプリケーション全体でformatを統一できるなら、このバリデーションを使っても良さそう
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	fmt.Println(err)
	if err != nil {
		return false
	}
	return true
}

func dateRangeValidation(fl validator.FieldLevel) bool {

	var date = fl.Field().String()
	var minDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	var maxDate = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	datetime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	if datetime.Before(minDate) || datetime.After(maxDate) {
		return false
	}
	return true
}

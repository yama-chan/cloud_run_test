package converter

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/taisukeyamashita/test/lib/errs"
)

// StructToMap structをmapに変換
func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	// Elem() : dataがinterface及びポインター型の場合、
	// 			interfaceに含まれるreflect.Value、またはポインターが参照しているreflect.Valueを返す
	// 			ただし、dataがinterface、ポインター型でもない場合はpanicになるので注意
	elem := reflect.ValueOf(data).Elem()
	// dataが構造体でない場合はエラーを返すようにする
	if elem.Kind() != reflect.Struct {
		return nil, errs.NewXerrorWithMessage(
			fmt.Sprintf("fail to StructToMap: %v is not struct", elem.Kind()),
		)
	}
	size := elem.NumField()

	for i := 0; i < size; i++ {
		// Field(): elemが構造体の場合、i番目のフィールドを取得
		// 			ただし、elemが構造体でない場合はpanicになるので注意
		// 			つまり、この関数ではelemが構造体でない場合はpanicになるので注意
		// Name : フィールド名

		// フィールド名を取得
		field := elem.Type().Field(i).Name
		// フィールドの値をinterface{}型で取得
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result, nil
}

// StructToJsonTagMap jsonタグを含むstructをmapに変換(パターン１)
func StructToJsonTagMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	// dataが構造体でない場合はエラーを返すようにする
	if elem.Kind() != reflect.Struct {
		return nil, errs.NewXerrorWithMessage(
			fmt.Sprintf("fail to StructToMap: %s is not struct", elem.Kind().String()),
		)
	}
	size := elem.NumField()

	for i := 0; i < size; i++ {
		// フィールドの　jsonタグに関連付けられた値を取得
		// 【例】`json:"name"`の場合は、nameを取得できる
		field := elem.Type().Field(i).Tag.Get("json")
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result, nil
}

// StructToJsonTagMap2 jsonタグを含むstructをmapに変換(パターン２)
func StructToJsonTagMap2(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	return result
}

/*********************************************************************

【使用例】
type A struct {
    ID int `json:"id"`
    Name string `json:"name"`
}

func main() {
    a := A{1, "keitaj"}
    fmt.Println(a)

    b := StructToMap(&a)
    fmt.Println(b)

    c := StructToJsonTagMap(&a)
    fmt.Println(c)

    d := StructToJsonTagMap2(&a)
    fmt.Println(d)
}


【実行結果】
{1 keitaj}
map[ID:1 Name:keitaj]
map[id:1 name:keitaj]
map[name:keitaj id:1]

*********************************************************************/

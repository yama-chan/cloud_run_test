package utils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/taisukeyamashita/test/lib/errs"
)

// StructToMap structをmapに変換
func StructToMap(data interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(data)
	// 構造体のポインタでない場合はエラー
	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Struct {
		return nil, errs.NewXerrorWithMessage(
			fmt.Sprintf("fail to StructToMap: [%v %v] is not [Ptr struct]", v.Kind(), v.Type().Elem().Kind()),
		)
	}

	result := make(map[string]interface{})
	// Elem() : dataがinterface及びポインター型の場合、
	// 			そのinterfaceに含まれるreflect.Value、またはポインターが参照しているreflect.Valueを返す
	// 			ただし、dataがinterface、ポインター型でもない場合はpanicになるので注意(nilの場合はnilを返す)
	//			故に、dataがポインター型でreflect.Value経由で値をセットしたい場合は　Elem()を使用することになる
	//			また、dataがinterface型で、そのKindが　ポインター型（ptr）の場合も上記と同様にElem()を使用して値を変更できる（以下、その例）
	//**************************************************
	//			v interface{}
	//			rv := reflect.ValueOf(v)
	// 			if rv.Kind() == reflect.Ptr {
	//     		rv = reflect.ValueOf(v).Elem()
	// 			}
	//			・・・
	//			rv.Set('セットする値')
	//			rv.SetInt('セットするInt')
	//**************************************************
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		// Field(): elemが構造体の場合、i番目のフィールドを取得
		// 			ただし、elemが構造体でない場合はpanicになるので注意
		// 			つまり、この関数ではelemが構造体でない場合はpanicになるので注意
		// Name : フィールド名

		// フィールド名を取得
		field := elem.Type().Field(i).Name
		// フィールドの値をinterface{}型で取得
		// MEMO: Elem()でreflect.Valueの指し示す先の値を取得していないとFieldメソッドを使った時にpanicになるので注意
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result, nil
}

// StructToJSONTagMap jsonタグを含むstructをmapに変換(パターン１)
func StructToJSONTagMap(data interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(data)
	// 構造体のポインタでない場合はエラー
	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Struct {
		return nil, errs.NewXerrorWithMessage(
			fmt.Sprintf("fail to StructToMap: [%v %v] is not [Ptr struct]", v.Kind(), v.Type().Elem().Kind()),
		)
	}
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
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

// StructToJSONTagMap2 jsonタグを含むstructをmapに変換(パターン２)
func StructToJSONTagMap2(data interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(data)
	// 構造体のポインタでない場合はエラー
	if v.Type().Kind() != reflect.Ptr || v.Type().Elem().Kind() != reflect.Struct {
		return nil, errs.NewXerrorWithMessage(
			fmt.Sprintf("fail to StructToMap: [%v %v] is not [Ptr struct]", v.Kind(), v.Type().Elem().Kind()),
		)
	}
	result := make(map[string]interface{})

	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	return result, nil
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

    c := StructToJSONTagMap(&a)
    fmt.Println(c)

    d := StructToJSONTagMap2(&a)
    fmt.Println(d)
}


【実行結果】
{1 keitaj}
map[ID:1 Name:keitaj]
map[id:1 name:keitaj]
map[name:keitaj id:1]

*********************************************************************/

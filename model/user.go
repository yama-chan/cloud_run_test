package model

import (
	"sync"
)

//Userの構造体
type User struct {
	UserInf UserInf
	mu      sync.Mutex //排他制御用　※一意のIDを生成するためユーザ同時登録を防ぐ
}

//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type UserInf struct {
	ID               string
	Fullname         string
	Sex              string
	Email            string
	Department       string
	Description      string
	RegisteredDate   string
	LastModifiedDate string
	IsDeleted        bool
	IconURL          string
	EncodedKey       string
}

//userのスライスを型宣言
//※一覧にて表示するため変数名を大文字で始めることでpublicな変数として扱う。
type Users []UserInf

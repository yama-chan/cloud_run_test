package utils

import (
	"fmt"

	"github.com/rs/xid"
)

//モデル用のユニークIDを生成
func CreateUniqueId() string {
	// ************************
	// xidについて
	// ************************
	// 詳しくはGitHubのREADMEの書かれていますが、その中から一部抜粋して紹介します。

	// binaryのformat
	// 全体で12bytesで、先頭から以下のように構成されています。

	// 4bytes: Unix timestamp (秒単位)
	// 3bytes: ホストの識別子
	// 2bytes: プロセスID
	// 3bytes: ランダムな値からスタートしたカウンタの値

	// 生成される文字列
	// 20文字のlower caseの英数字。([0-9a-v]{20})
	// 例: b8hpcg8hv3amvi9dol0g

	// idを生成
	guid := xid.New()
	fmt.Println("created unique id :" + guid.String())

	// binaryの各partの情報を取得（参考に）
	machine := guid.Machine()
	pid := guid.Pid()
	time := guid.Time()
	counter := guid.Counter()
	fmt.Printf("machine: %v,\n pid: %v,\n time: %v,\n counter: %v\n", machine, pid, time, counter)

	return guid.String()
}

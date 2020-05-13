package route

import "github.com/taisukeyamashita/test/lib"

// Route ルート
type Route struct {
	Gets    []lib.EndPoint
	Posts   []lib.EndPoint
	Puts    []lib.EndPoint
	Deletes []lib.EndPoint
	Patches []lib.EndPoint
}

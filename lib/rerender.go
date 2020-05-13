package lib

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// JSONResponseContent JSONレスポンスコンテンツ
type JSONResponseContent interface{}

// FileResponseContent ファイルレスポンスコンテンツ
type FileResponseContent struct {
	Filename    string
	Data        []byte
	ContentType string
}

// ResponseRenderer レスポンスを吐き出す
type ResponseRenderer interface {
	Render(echo.Context) error
}

// JSONResponseRenderer JSONをレスポンスとして吐き出す
type JSONResponseRenderer struct {
	Status  int
	Content JSONResponseContent
}

// FileResponseRenderer ファイルをレスポンスとして吐き出す
type FileResponseRenderer struct {
	Status  int
	Content FileResponseContent
}

// Render implements ResponseRenderer
func (j JSONResponseRenderer) Render(ctx echo.Context) error {
	return ctx.JSON(j.Status, j.Content)
}

// Render implements ResponseRenderer
func (f FileResponseRenderer) Render(ctx echo.Context) error {
	// ブラウザで読み込むのではなく、ファイルを直接ダウンロードさせる
	disposition := "attachment"
	if f.Content.Filename != "" {
		disposition += fmt.Sprintf(`; filename="*=utf-8'';%s"`, url.PathEscape(f.Content.Filename))
	}
	ctx.Response().Header().Add("Content-Disposition", disposition)
	ctx.Response().Header().Add("X-FILE-NAME", url.PathEscape(f.Content.Filename))
	return ctx.Blob(f.Status, f.Content.ContentType, f.Content.Data)
}

// RenderFileOK 成功が200のファイルコンテンツをレスポンスに変換
func RenderFileOK(content FileResponseContent, err error) (ResponseRenderer, error) {
	return renderFileResponse(http.StatusOK, content, err)
}

// RenderOK 成功が200のコンテンツをレスポンスに変換
func RenderOK(content JSONResponseContent, err error) (ResponseRenderer, error) {
	return renderResponse(http.StatusOK, content, err)
}

// RenderCreated 成功が201のコンテンツをレスポンスに変換
func RenderCreated(content JSONResponseContent, err error) (ResponseRenderer, error) {
	return renderResponse(http.StatusCreated, content, err)
}

// RenderNoContent 成功が204のコンテンツをレスポンスに変換
func RenderNoContent(err error) (ResponseRenderer, error) {
	return renderResponse(http.StatusNoContent, nil, err)
}

func renderFileResponse(successStatus int, content FileResponseContent, err error) (ResponseRenderer, error) {
	if err != nil {
		return nil, err
	}
	return FileResponseRenderer{Status: successStatus, Content: content}, nil
}

func renderResponse(succeedStatus int, content JSONResponseContent, err error) (ResponseRenderer, error) {
	if err != nil {
		return nil, err
	}
	return JSONResponseRenderer{Status: succeedStatus, Content: content}, nil
}

// EndPointHandler 独自のResponseRendererを返す関数
// ここでカスタムContextを定義して引数にしても可。ルーティング時に実行される関数として定義する。
type EndPointHandler func(echo.Context) (ResponseRenderer, error)

// EndPoint エンドポイント
type EndPoint struct {
	Path    string
	Handler EndPointHandler
}

// NewEndPoint 新しいエンドポイントを作成
func NewEndPoint(path string, handler EndPointHandler) EndPoint {
	return EndPoint{Path: path, Handler: handler}
}

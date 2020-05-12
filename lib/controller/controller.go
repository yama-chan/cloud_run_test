package controller

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
	"github.com/taisukeyamashita/test/lib/errs"
	"github.com/taisukeyamashita/test/lib/server"
)

// Controller コントローラ
type Controller interface {
	RegistControllers(mux *http.ServeMux)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// ControllerBase 既定コントローラ
type ControllerBase struct {
	Echo     *echo.Echo
	Provider server.Provider
}
// ControllerBase implement http.Handler
var _ http.Handler = ControllerBase{}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (controller ControllerBase) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Echo.ServeHTTP(w, r)
}

// NewController Controllerを作成
func NewController(provider server.Provider) ControllerBase {
	e := echo.New()
	base := ControllerBase{Engine: e, Provider: provider}

	// BEFORE middlewares
	e.Use(
		// HandlerFuncの実行前処理
		base.withContextGen(),
		// base.withCustomRecover(),
	)

	// AFTER middleware
	// e.Use(
	// 	base.withProviderFinalizer(),
	// )

	return base
}

// AddRoutes Route登録
func (controller ControllerBase) AddRoutes(group *echo.Group, router interface{}) {
	reflectedRouter := reflect.ValueOf(router).Elem()
	providerValue := reflect.ValueOf(controller.Provider)
	for index := 0; index < reflectedRouter.NumMethod(); index++ {
		method := reflectedRouter.Method(index)

		result := method.Call([]reflect.Value{providerValue})
		route := result[0].Interface().(Route)
		controller.addRoute(group, route)
	}
}

func (controller ControllerBase) addRoute(group *echo.Group, route Route) {
	controller.addEndPoints("GET", group, route.Gets)
	controller.addEndPoints("POST", group, route.Posts)
	controller.addEndPoints("PUT", group, route.Puts)
	controller.addEndPoints("DELETE", group, route.Deletes)
	controller.addEndPoints("PATCH", group, route.Patches)
}

func (controller ControllerBase) addEndPoints(method string, group *echo.Group, endPoints []EndPoint) {
	for _, endPoint := range endPoints {
		group.Add(method, endPoint.Path, controller.endPointHandlerToEchoHandler(endPoint.Handler))
	}
}

func (controller ControllerBase) endPointHandlerToEchoHandler(handler endPointHandler) echo.HandlerFunc {
	return func(context echo.Context) error {
		// l := controller.Provider.Logger(context.Request().Context())
		r, e := handler(newContext(context, l))
		if e != nil {
			return errs.WrapXerror(e)
		}
		return r.Render(context)
	}
}

// ***********************************************************************
// https://medium.com/veltra-engineering/echo-middleware-in-golang-90e1d301eb27

//	Middlewareの実行順序

// 	middleware-Pre  : before
// 	middleware-Use-1: before
// 	middleware-Use-2: before
// 	middleware-Group: before
// 	middleware-Route: before
// 	logic: main
// 	logic: defer
// 	middleware-Route: after
// 	middleware-Route: defer
// 	middleware-Group: after
// 	middleware-Group: defer
// 	middleware-Use-2: after
// 	middleware-Use-2: defer
// 	middleware-Use-1: after
// 	middleware-Use-1: defer
// 	middleware-Pre  : after
// 	middleware-Pre  : defer

//	★ 'Pre'→'Use'→'Group'→'Route'の順
//	★ 'Use'で設定された2つについては、先に設定したものから実行されている
//	★ 'defer'が実行されるタイミングは当該Middlewareの事後処理('after')直後

// ***********************************************************************

func (controller ControllerBase) withContextGen() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		//defer内部で発生したerrorを処理するのには名前付き返り値を利用する。
		return func(c echo.Context) (err error) {
			// echoのhttp.Requestのをcontext.Contextをアプリケーション用のcontext.Contextに成形
			ctx := controller.Provider.Context(c.Request())
			// Middlewareの事後処理(after)直後にdeferが実行されるので、ここでFinalize
			defer func() {
				ferr := controller.Provider.Finalize(ctx)
				if ferr == nil {
					return
				}
				// 名前付き返り値であるerrに上書き
				err = fmt.Errorf("Failed to close: %v, the original error was %v", cerr, err)
			}() // カッコ'()'で実行（※関数型変数fをf()で実行するイメージ）

			// WithContext : 引数のctxに書き換えた*http.Requestの'コピー'を新規で生成する
			// SetRequest : 引数の*http.Requestをecho.Contextにセットする
			c.SetRequest(c.Request().WithContext(WithContext))
			// ↑ BEFORE
			// この場合、HandlerFuncが実行されてReturnとなる
			return next(c) // HandlerFunc : func(Context) error
			// ↓ AFTER
		}
	}
}

func (controller ControllerBase) withProviderFinalizer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ↑ BEFORE
			err := next(c) // HandlerFunc : func(Context) error
			// この場合、AFTERの処理は実行され、エラーを返す
			// ↓ AFTER
			if err != nil {
				return err
			}
			finalizeError := controller.Provider.Finalize(c.Request().Context())
			if finalizeError != nil {
				return finalizeError
			}
			return nil
		}
	}
}

func (controller ControllerBase) withProviderClient() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			controller.Provider.

			// ↑ BEFORE
			err := next(c) // HandlerFunc : func(Context) error
			// この場合、AFTERの処理は実行され、エラーを返す
			// ↓ AFTER
			if err != nil {
				return err
			}
			finalizeError := controller.Provider.Finalize(c.Request().Context())
			if finalizeError != nil {
				return finalizeError
			}
			return nil
		}
	}
}

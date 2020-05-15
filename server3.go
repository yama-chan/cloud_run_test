package main

//**************************************************************************
// TODO: golang にて　TCP　でソケット通信できるサーバー実装を検討した
//**************************************************************************
import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	_ "syscall"
)

func addRoutes() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world.")
	})
	return
}

// onceCloseListener wraps a net.Listener, protecting it from multiple Close calls.
// Closeは一度だけを保証する
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }

// TODO: 不適切ならコメントに戻す
// func listenServer(ctx context.Context, addr string) (l net.Listener, ch chan error) {
func listenServer(ctx context.Context, addr string) (l net.Listener, ch chan error) {
	ch = make(chan error)

	// 「http.ListenAndServe」だと”tcp”固定となるみたい
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("⇨ http server started on %s\n", listener.Addr())

	// Closeは一度だけ呼ばれるようにラップする
	l = &onceCloseListener{Listener: listener}
	defer l.Close()
	// TODO: リスナーでの実装を検討、まだソケット処理を実装出来ていない...
	// TODO: Echoの実装のようにリスナーやエラーを何かしらの構造体でラップしてもいいのかな？
	// con := l.Accept()

	// *http.Serverを作成
	server := &http.Server{
		Addr: addr,
		// Handler: mux,
	}
	go func() {
		mux := addRoutes()
		server.Handler = mux
		// TODO: (*http.Server).Serve　と　http.Serve(l net.Listener, handler http.Handler)　の違いはあるのか要確認

		// Serveはリスナー（'l'）で着信HTTP接続を受け入れ、それぞれに新しいサービスゴルーチンを作成します。
		// サービスゴルーチンは要求を読み取り、次にハンドラーを呼び出してそれらに応答します。
		// ch <- http.Serve(l, mux)
		// Serveはリスナー（'l'）で着信接続を受け入れ、それぞれに新しいサービスゴルーチンを作成します。
		// サービスgoroutineは要求を読み取り、srv.Handlerを呼び出して要求に応答します。
		ch <- server.Serve(l) // 名前付き返り値(ch)に上書き
	}()
	return
}

// TODO: flagを理解したら学習用途で盛り込みたい
// var (
// 	proto = flag.String("proto", "tcp", "proto(tcp|tcp4|tcp6)")
// 	addr  = flag.String("addr", "0.0.0.0", "listen addres")
// 	port  = flag.Int("port", 9999, "listen port")
// )

func run3(stdout, stderr io.Writer, args []string) int {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	listener, errCh := listenServer(ctx, ":8080")
	fmt.Println("Server started at", listener.Addr())

	// シグナルハンドリング (Ctrl + C)
	sigCh := make(chan os.Signal)
	//何かしらの処理を中断するシグナルがあった場合は、listenerをClsoeする
	// signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sigCh, os.Interrupt)
	// go func() {
	// 	log.Println(<-sig)
	// 	listener.Close()
	// }()

	// TODO: チャネルにてエラー取得時のエラーハンドリングを理解する
	log.Println(<-errCh)
	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("info: server encounters an error: %v", err)
		} else {
			log.Print("info: server has stopped unexpectedly")
		}
		return 1
	case sig := <-sigCh:
		cancel()
		lerr := listener.Close() // リスナーを明示的にClose()
		if lerr != nil {
			return 1
		}
		if err := <-errCh; errors.Is(err, context.Canceled) {
			return 0
		} else if err != nil {
			log.Printf("info: failed to stop server by receiving a signal %s: %v", sig, err)
			return 1
		} else {
			return 0
		}
	case <-ctx.Done():
		return 0
	}
}

func main3() {
	// 明示的に引数を書く、コードで返してExitすることでテストを書きやすくなる
	os.Exit(run3(os.Stdout, os.Stderr, os.Args))
}

package main

//**************************************
// TODO: golang にて　TCP　でソケット通信する際のサーバー実装を検討
//**************************************
import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	_ "syscall"
	"time"
)

func addRoutes() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world.")
	})
	return
}

// onceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
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

func listenServer(ctx context.Context, addr string) (net.Listener, error) {
	ch := make(chan error)

	// 「http.ListenAndServe」だと”tcp”固定となるみたい
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("⇨ http server started on %s\n", listener.Addr())

	// Closeは一度だけ呼ばれるようにラップする
	l := &onceCloseListener{Listener: listener}
	defer l.Close()

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
	select {
	case err := <-ch:
		return l, err
	case <-ctx.Done():
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("warn: failed to stop the server gracefully: %v", err)
		}
		return l, ctx.Err()
	}
	return l, nil
}

// TODO: flagを理解したら盛り込みたい
// var (
// 	proto = flag.String("proto", "tcp", "proto(tcp|tcp4|tcp6)")
// 	addr  = flag.String("addr", "0.0.0.0", "listen addres")
// 	port  = flag.Int("port", 9999, "listen port")
// )

func main3() {
	ctx := context.Background()
	listener, ch := listenServer(ctx, ":8080")
	fmt.Println("Server started at", listener.Addr())

	// シグナルハンドリング (Ctrl + C)
	sig := make(chan os.Signal)

	//何かしらの処理を中断するシグナルがあった場合は、listenerをClsoeする
	// signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, os.Interrupt)
	go func() {
		log.Println(<-sig)
		listener.Close()
	}()

	// TODO: チャネルにてエラー取得時のエラーハンドリングを追加する
	log.Println(<-ch)
}

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

var logger = log.Default()

func logHttp(req *http.Request) {
	if dump, err := httputil.DumpRequest(req, false); err == nil {
		logger.Printf("%s %q\n", req.RemoteAddr, dump)
	} else {
		logger.Println(err)
	}
}

type handler struct{}

func (*handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	go logHttp(req)

	url := *req.URL
	url.Scheme = "https"
	url.Host = req.Host

	http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
}

func listenAndServe(addr string) error {
	return (&http.Server{
		Addr:    addr,
		Handler: &handler{},
		ConnState: func(conn net.Conn, state http.ConnState) {
			go func() {
				if state == http.StateActive || state == http.StateIdle {
					return
				}
				logger.Printf("%s %v", conn.RemoteAddr(), state)
			}()
		},
		BaseContext: func(lst net.Listener) context.Context {
			go logger.Printf("%s listen", lst.Addr())
			return context.Background()
		},
	}).ListenAndServe()
}

func main() {
	var addr string
	argc := len(os.Args)
	switch {
	case argc == 2:
		addr = os.Args[1]
	case argc == 3 && os.Args[1] == "/upgrade":
		addr = os.Args[2]
	default:
		log.Fatalf("usage: %s ADDR\n", os.Args[0])
	}
	log.Fatalln(listenAndServe(addr))
}

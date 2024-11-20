package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
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

func main() {
	srv := &http.Server{
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
	}
	err := srv.ListenAndServe()
	log.Fatalln(err)
}

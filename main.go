package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type handler struct{}

var logger = log.Default()

func dump(req *http.Request) {
	if dump, err := httputil.DumpRequest(req, false); err == nil {
		logger.Printf("%q\n", dump)
	} else {
		logger.Println(err)
	}
}

func (*handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dump(req)

	url := *req.URL
	url.Scheme = "https"
	url.Host = req.Host
	req.Header.Set("content-type", "")
	http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
}

func main() {
	log.Fatalln(http.ListenAndServe(":80", &handler{}))
}

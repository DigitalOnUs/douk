package main

import (
	"net/http"

	"flag"

	"goji.io"
	"goji.io/pat"
)

// APIs
var (
	address string
)

func init() {
	flag.StringVar(&address, "listen", "localhost:7001", "listen address")
}

func main() {
	mux := goji.NewMux()
	flag.Parse()

	defaultPath := "./statics"

	mux.Handle(pat.Get("/*"), http.FileServer(http.Dir(defaultPath)))
	http.ListenAndServe(address, mux)
}

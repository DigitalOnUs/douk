//go:generate go run github.com/UnnoTed/fileb0x b0x.json
package main

import (
	"flag"
	"net/http"

	"github.com/DigitalOnUs/douk/api"
	"github.com/DigitalOnUs/douk/static"
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

	// For debugging replace the static.HTTP to defaultPath
	// and _ comment the import for statics pre-compiled in the binary
	mux.HandleFunc(pat.Post("/api/consulize"), api.Consulize)
	mux.Handle(pat.Get("/*"), http.FileServer(static.HTTP))
	//mux.Handle(pat.Get("/*"), http.FileServer(http.Dir("./public")))

	http.ListenAndServe(address, mux)
}

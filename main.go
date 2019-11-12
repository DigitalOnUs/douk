//go:generate go run github.com/UnnoTed/fileb0x b0x.json
package main

import (
	"flag"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/DigitalOnUs/douk/static"
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
	//defaultPath := "./public"

	mux.Handle(pat.Get("/*"), http.FileServer(static.HTTP))
	http.ListenAndServe(address, mux)
}
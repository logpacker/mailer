//go:generate swagger generate spec -o ./swagger.json
// Mailer API
//
//     Schemes: http
//     Host: 127.0.0.1:6100
//     BasePath: /v1
//     Version: 0.0.1
//     Contact: pltvs@logpacker.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"github.com/logpacker/mailer/pkg/api"
	"log"
	"net/http"
	"os"
)

// Version var
var Version string

func main() {
	version := flag.Bool("v", false, "Build version (git revision)")
	help := flag.Bool("h", false, "Usage & Help")
	apiKey := flag.String("a", "", "Ser secret api_key")
	flag.String("s", "localhost:25", "SMTP address")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	r := api.NewRouter(*apiKey)
	log.Println(http.ListenAndServe(":6100", context.ClearHandler(r)))
}

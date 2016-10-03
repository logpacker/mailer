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

// APIKey var
var APIKey string

func main() {
	version := flag.Bool("v", false, "cmd version")
	apiKey := flag.String("a", "", "authorized api key")
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *apiKey != "" {
		APIKey = *apiKey
	}

	r := api.NewRouter(APIKey)
	log.Println(http.ListenAndServe(":6100", context.ClearHandler(r)))
}

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

//go:generate swagger generate spec -o ./swagger.json
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"github.com/logpacker/mailer/pkg/api"
	"github.com/logpacker/mailer/pkg/shared"
	"net/http"
	"os"
)

// Version var
var Version string

func main() {
	version := flag.Bool("v", false, "Build version (git revision)")
	help := flag.Bool("h", false, "Usage & Help")
	apiKey := flag.String("a", "", "Set secret api_key. If empty API will be accessible without token")
	p := flag.String("p", "6100", "API port to bind")
	b := flag.String("b", "127.0.0.1:11300", "Beanstalkd connection string")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	conf := new(shared.MailerConfig)
	conf.BeanstalkdAddr = *b

	r := api.NewRouter(*apiKey, conf)
	err := http.ListenAndServe(fmt.Sprintf(":%s", *p), context.ClearHandler(r))
	shared.LogErr(err)
}

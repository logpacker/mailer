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
	"github.com/logpacker/mailer/pkg/conf"
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
	db := flag.String("db", "root@tcp(127.0.0.1:3306)/mailer", "MySQL database connection string")
	*db += "?charset=utf8&parseTime=true"
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	conf := new(conf.MailerConfig)
	conf.MySQLAddr = *db

	r := api.NewRouter(*apiKey, conf)
	err := http.ListenAndServe(fmt.Sprintf(":%s", *p), context.ClearHandler(r))
	shared.LogErr(err)
}

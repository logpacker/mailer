package main

import (
	"flag"
	"fmt"
	"github.com/logpacker/mailer/pkg/conf"
	"github.com/logpacker/mailer/pkg/daemon"
	"os"
)

// Version var
var Version string

func main() {
	version := flag.Bool("v", false, "Build version (git revision)")
	help := flag.Bool("h", false, "Usage & Help")
	smtp := flag.String("s", "localhost:25", "SMTP address")
	db := flag.String("db", "root@tcp(127.0.0.1:3306)/mailer", "MySQL database connection string")
	b := flag.String("b", "127.0.0.1:11300", "Beanstalkd connection string")
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
	conf.SMTPAddr = *smtp
	conf.MySQLAddr = *db
	conf.BeanstalkdAddr = *b

	daemon.StartConsumer(conf)
}

#### Dependencies

 * MySQL `CREATE DATABASE mailer CHARACTER SET utf8 COLLATE utf8_general_ci;` (import `db/schema.sql`)
 * SMTP
 * [Beanstalkd](http://kr.github.io/beanstalkd/)
 * [glide](https://github.com/Masterminds/glide)

#### Development

 * [go-swagger](https://github.com/go-swagger/go-swagger)
 * Genereate swagger.json - `go generate cmd/api/main.go`
 * Run tests - `go test ./pkg/... -v -cover`

#### Build & Deploy

```bash
cd $GOPATH/src
mkdir -p github.com/logpacker
cd github.com/logpacker
git clone git@github.com:logpacker/mailer.git
cd mailer
glide i
go build -ldflags "-X main.Version=$(git rev-parse HEAD)" -o mailer_api cmd/api/main.go
go build -ldflags "-X main.Version=$(git rev-parse HEAD)" -o mailer_daemon cmd/daemon/main.go
```

#### Usage and Flags

```bash
./mailer_api -h
./mailer_daemon -h
```

#### How does it work

 * `Sender` sends request to `Mailer API` to save mail into the queue
 * `Mailer Daemon` sends emails from queue to `Recepient` via SMTP and updates DB

Mail statuses:

 * Pending
 * Processing
 * Sent
 * Failed to Send
 * Opened

#### CURL Example

```bash
curl -H "Content-Type: application/json" \
-XPOST localhost:80/v1/send \
-d '{"from": {"email": "mailer@logpacker.com", "name": "LogPacker"}, "to": {"email": "alexander.plutov@gmail.com"}, "subject": "Verify your email address", "Body": "Thank you for the registration.<br/>Now please confirm it.", "url_unsubscribe": "https://logpacker.com"}'
```

#### Go Client example

```go
package main

import (
	"github.com/logpacker/mailer/pkg/client"
)

func main() {
	// Error handling skipped for better readability
	c, _ := client.New(client.Config{
		URL:    "http://127.0.0.1:80",
		APIKey: "secret",
	})

	c.SendEmail(client.Email{
		From: &client.Address{
			Email: "mailer@logpacker.com",
			Name:  "LogPacker",
		},
		To: &client.Address{
			Email: "alexander.plutov@gmail.com",
		},
		Subject:        "Verify your email address",
		Body:           "Thank you for the registration.<br/>Now please confirm it.",
		URLUnsubscribe: "https://logpacker.com",
	})
}

```

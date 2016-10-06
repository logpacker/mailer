#### Dependencies

 * MySQL
 * SMTP
 * Beanstalkd

#### Deployment dependencies:

 * [goose](https://bitbucket.org/liamstask/goose/). Create a DB first: `CREATE DATABASE mailer CHARACTER SET utf8 COLLATE utf8_general_ci;`
 * [glide](https://github.com/Masterminds/glide)

#### Development

 * [go-swagger](https://github.com/go-swagger/go-swagger)
 * `go generate cmd/api/main.go`
 * `go test ./pkg/... -v -cover`

#### Build & Deploy

```bash
cd $GOPATH/src
mkdir -p github.com/logpacker
cd github.com/logpacker
git clone git@github.com:logpacker/mailer.git
cd mailer
glide i
goose --env=live up
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
-XPOST localhost:6100/v1/send \
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
		URL:    "http://127.0.0.1:6100",
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

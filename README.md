#### Dependencies

 * MySQL
 * SMTP
 * Beanstalkd

#### Deployment dependencies:

 * [goose](https://bitbucket.org/liamstask/goose/). Create a DB first: `CREATE DATABASE mailer CHARACTER SET utf8 COLLATE utf8_general_ci;`
 * [glide](https://github.com/Masterminds/glide)
 * `go test ./pkg/... -v -cover`

#### Development

 * [go-swagger](https://github.com/go-swagger/go-swagger)

#### Build

```bash
go build -ldflags "-X main.Version=$(git rev-parse HEAD)" -o mailer_api cmd/api/main.go
go build -ldflags "-X main.Version=$(git rev-parse HEAD)" -o mailer_daemon cmd/daemon/main.go
# Generate swagger.json
go generate cmd/api/main.go
```

#### Usage and Flags

```bash
./mailer_api -h
./mailer_daemon -h
```

#### How does it work

 * `Sender` sends request to `Mailer API` to save mail into the queue and DB
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
-d '{"from": {"email": "mailer@logpacker.com", "name": "LogPacker"}, "to": {"email": "alexander.plutov@gmail.com"}, "subject": "Verify your email address", "Body": "<b>Thank you for the registration. Now please confirm it.</b>", "url_unsubscribe": "http://logpacker.com/unsubscribe"}'
```

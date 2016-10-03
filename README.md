#### Dependencies

 * MySQL

#### Deployment dependencies:

 * [goose](https://bitbucket.org/liamstask/goose/)
 * [glide](https://github.com/Masterminds/glide)

#### Development

 * [go-swagger](https://github.com/go-swagger/go-swagger)

#### Build

```bash
go build -ldflags "-X main.Version=$(git rev-parse HEAD) -X main.APIKey=secret" -o mailer_api cmd/api/main.go
go generate cmd/api/main.go
```

#### Usage

 - `./mailer_api -v` - shows build version
 - `./mailer_api -a secret` - starts API with 'secret' api_key

#### How does it work

 * `Sender` sends request to `Mailer API` to save mail into the queue
 * `Mailer Daemon` sends emails from queue to `Recepient` via SMTP

Mail statuses:

 * Pending
 * Sent
 * Failed to Send
 * Opened

#### Dependencies

 * MySQL
 * SMTP

#### Deployment dependencies:

 * [goose](https://bitbucket.org/liamstask/goose/). Create a DB first: `CREATE DATABASE mailer CHARACTER SET utf8 COLLATE utf8_general_ci;`
 * [glide](https://github.com/Masterminds/glide)

#### Development

 * [go-swagger](https://github.com/go-swagger/go-swagger)

#### Build

```bash
go build -ldflags "-X main.Version=$(git rev-parse HEAD)" -o mailer_api cmd/api/main.go
go generate cmd/api/main.go
```

#### Usage and Flags

```bash
./mailer_api -h

Usage of mailer_api:
  -a string
    	Ser secret api_key
  -h	Usage & Help
  -s string
    	SMTP address (default "localhost:25")
  -v	Build version (git revision)
```

#### How does it work

 * `Sender` sends request to `Mailer API` to save mail into the queue
 * `Mailer Daemon` sends emails from queue to `Recepient` via SMTP

Mail statuses:

 * Pending
 * Sent
 * Failed to Send
 * Opened

package daemon

import (
	"github.com/logpacker/mailer/pkg/conf"
	"github.com/logpacker/mailer/pkg/db"
	"github.com/logpacker/mailer/pkg/queue"
	"github.com/logpacker/mailer/pkg/shared"
)

var (
	dbClient    *db.MySQLClient
	queueClient *queue.BeanstalkdClient
	smtpClient  *SMTPClient
)

// StartConsumer func
func StartConsumer(conf *conf.MailerConfig) {
	dbClient = new(db.MySQLClient)
	dbErr := dbClient.Init(conf.MySQLAddr)
	if dbErr != nil {
		panic(dbErr)
	}

	queueClient = new(queue.BeanstalkdClient)
	queueErr := queueClient.Init(conf.BeanstalkdAddr)
	if queueErr != nil {
		panic(queueErr)
	}

	smtpClient = new(SMTPClient)
	smtpClient.Init(conf.SMTPAddr)

	queueClient.ReceiveEmails(SendEmail)
}

// SendEmail used as a callback
func SendEmail(email *shared.Email) {
	var err error
	err = dbClient.UpdateStatus(email, db.StatusProcessing)
	shared.LogErr(err)
	if err != nil {
		return
	}

	smtpEmail := BuildSMTPEmail(email)
	err = smtpClient.Send(smtpEmail)
	shared.LogErr(err)
	resultStatus := db.StatusSent
	if err != nil {
		resultStatus = db.StatusFailedToSend
	}

	err = dbClient.UpdateStatus(email, resultStatus)
	shared.LogErr(err)
}

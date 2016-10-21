package daemon

import (
	"github.com/logpacker/mailer/pkg/db"
	"github.com/logpacker/mailer/pkg/queue"
	"github.com/logpacker/mailer/pkg/shared"
	"sync"
)

var (
	dbClient    *db.MySQLClient
	queueClient *queue.BeanstalkdClient
	smtpClient  *SMTPClient
)

// StartConsumer func
func StartConsumer(conf *shared.MailerConfig) {
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
	smtpClient.Conf = conf
	smtpClient.Init(conf.SMTPAddr)

	var wg sync.WaitGroup
	for i := 0; i < conf.ConsumersCount; i++ {
		wg.Add(1)
		go queueClient.ReceiveEmails(SendEmail)
		wg.Add(1)
		go queueClient.ReceiveOpenEmails(OpenEmail)
	}
	wg.Wait()
}

// SendEmail used as a callback
func SendEmail(email *shared.Email) {
	var err error

	err = dbClient.SaveEmail(email)
	if err != nil {
		shared.Logf("Unable to save email into the DB. Detail: %s", err.Error())
		return
	}
	shared.Logf("Email saved. ID: %d, To: %s", email.ID, email.To.Email)

	err = dbClient.UpdateStatus(email, db.StatusProcessing)
	shared.LogErr(err)
	if err != nil {
		return
	}

	smtpEmail := BuildSMTPEmail(email, smtpClient.Conf)
	err = smtpClient.Send(smtpEmail)
	shared.LogErr(err)
	resultStatus := db.StatusSent
	if err != nil {
		resultStatus = db.StatusFailedToSend
	}

	err = dbClient.UpdateStatus(email, resultStatus)
	shared.LogErr(err)

	err = dbClient.UpdateSentAt(email)
	shared.LogErr(err)
}

// OpenEmail used as a callback
func OpenEmail(openEmail *shared.OpenEmail) {
	var err error
	email := &shared.Email{
		ID: openEmail.ID,
	}

	shared.Logf("Email opened. ID: %d", email.ID)
	err = dbClient.UpdateStatus(email, db.StatusOpened)
	shared.LogErr(err)

	err = dbClient.UpdateOpenedAt(email)
	shared.LogErr(err)
}

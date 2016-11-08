package queue

import (
	"encoding/json"
	"github.com/kr/beanstalk"
	"github.com/logpacker/mailer/pkg/shared"
	"time"
)

var (
	queueEmails = "emails"
	queueOpen   = "open"
)

// SendCallback func type
type SendCallback func(email *shared.Email)

// OpenCallback func type
type OpenCallback func(email *shared.OpenEmail)

// BeanstalkdClient struct
type BeanstalkdClient struct {
	Conn *beanstalk.Conn
}

// Init func
func (s *BeanstalkdClient) Init(addr string) error {
	var err error
	s.Conn, err = beanstalk.Dial("tcp", addr)
	return err
}

// SendEmailJob func
func (s *BeanstalkdClient) SendEmailJob(email *shared.Email) error {
	var (
		id   uint64
		err  error
		data []byte
	)

	data, err = json.Marshal(email)
	if err != nil {
		return err
	}

	tube := &beanstalk.Tube{
		Conn: s.Conn,
		Name: queueEmails,
	}
	id, err = tube.Put(data, 1, 0, time.Second)
	if err == nil {
		shared.Logf("Sent Send Job Id: %d", id)
	} else {
		shared.LogErr(err)
	}

	return err
}

// SendOpenJob func
func (s *BeanstalkdClient) SendOpenJob(openEmail *shared.OpenEmail) error {
	var (
		id   uint64
		err  error
		data []byte
	)

	data, err = json.Marshal(openEmail)
	if err != nil {
		return err
	}

	tube := &beanstalk.Tube{
		Conn: s.Conn,
		Name: queueOpen,
	}
	id, err = tube.Put(data, 1, 0, time.Second)
	if err == nil {
		shared.Logf("Sent Open Job Id: %d", id)
	} else {
		shared.LogErr(err)
	}

	return err
}

// ReceiveEmails func
func (s *BeanstalkdClient) ReceiveEmails(callback SendCallback) {
	tube := beanstalk.NewTubeSet(s.Conn, queueEmails)

	for {
		id, job, err := tube.Reserve(time.Second)
		if err == nil {
			processErr := s.ProcessSendEmailJob(id, job, callback)
			if processErr != nil {
				shared.Logf("Unable to process Send job id=%d. Details: %s", id, processErr.Error())
			}
		}
	}
}

// ReceiveOpenEmails func
func (s *BeanstalkdClient) ReceiveOpenEmails(callback OpenCallback) {
	tube := beanstalk.NewTubeSet(s.Conn, queueOpen)

	for {
		id, job, err := tube.Reserve(time.Second)
		if err == nil {
			processErr := s.ProcessOpenEmailJob(id, job, callback)
			if processErr != nil {
				shared.Logf("Unable to process Open job id=%d. Details: %s", id, processErr.Error())
			}
		}
	}
}

// ProcessSendEmailJob func
func (s *BeanstalkdClient) ProcessSendEmailJob(id uint64, job []byte, callback SendCallback) error {
	email := new(shared.Email)
	err := json.Unmarshal(job, email)

	shared.Logf("Processing Send Job Id: %d", id)
	s.DeleteJob(id)

	callback(email)

	return err
}

// ProcessOpenEmailJob func
func (s *BeanstalkdClient) ProcessOpenEmailJob(id uint64, job []byte, callback OpenCallback) error {
	email := new(shared.OpenEmail)
	err := json.Unmarshal(job, email)

	shared.Logf("Processing Open Job Id: %d", id)
	s.DeleteJob(id)

	callback(email)

	return err
}

// DeleteJob func
func (s *BeanstalkdClient) DeleteJob(id uint64) error {
	err := s.Conn.Delete(id)
	if err != nil {
		shared.Logf("Unable to delete job: id=%d. Details: %s", id, err.Error())
	}

	return err
}

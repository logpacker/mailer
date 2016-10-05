package queue

import (
	"encoding/json"
	"github.com/kr/beanstalk"
	"github.com/logpacker/mailer/pkg/shared"
	"time"
)

var (
	queueEmails = "emails"
)

// Callback func type
type Callback func(email *shared.Email)

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
	id, err = tube.Put(data, 1, 0, time.Minute)
	if err == nil {
		shared.Logf("Sent Job Id: %d", id)
	}

	return err
}

// ReceiveEmails func
func (s *BeanstalkdClient) ReceiveEmails(callback Callback) {
	tube := beanstalk.NewTubeSet(s.Conn, queueEmails)

	for {
		id, job, err := tube.Reserve(time.Minute)
		if err == nil {
			processErr := s.ProcessJob(id, job, callback)
			if processErr != nil {
				shared.Logf("Unable to process job id=%d. Details: %s", id, processErr.Error())
			}
		}
	}
}

// ProcessJob func
func (s *BeanstalkdClient) ProcessJob(id uint64, job []byte, callback Callback) error {
	shared.Logf("Processing Job Id: %d", id)

	email := new(shared.Email)
	err := json.Unmarshal(job, email)

	s.DeleteJob(id)

	go callback(email)

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

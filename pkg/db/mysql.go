package db

import (
	"database/sql"
	// Commented for lint
	_ "github.com/go-sql-driver/mysql"
	"github.com/logpacker/mailer/pkg/shared"
	"time"
)

var (
	// StatusPending var
	StatusPending = "Pending"

	// StatusProcessing var
	StatusProcessing = "Processing"

	// StatusSent var
	StatusSent = "Sent"

	// StatusFailedToSend var
	StatusFailedToSend = "Failed to Sent"

	// StatusOpened var
	StatusOpened = "Opened"
)

// MySQLClient struct
type MySQLClient struct {
	Conn  *sql.DB
	Stmts map[string]*sql.Stmt
}

// Init func
func (s *MySQLClient) Init(addr string) error {
	var err error
	s.Conn, err = sql.Open("mysql", addr)
	s.Stmts = make(map[string]*sql.Stmt)

	if err == nil {
		s.Stmts["get_address_id"], _ = s.Conn.Prepare("SELECT id FROM address WHERE email = ? AND is_sender = ?")
		s.Stmts["insert_address"], _ = s.Conn.Prepare("INSERT INTO address (email, name, is_sender) VALUES (?, ?, ?)")
		s.Stmts["insert_email"], _ = s.Conn.Prepare("INSERT INTO email (`from`, `to`, subject, body, url_unsubscribe) VALUES (?, ?, ?, ?, ?)")
		s.Stmts["update_email_status"], _ = s.Conn.Prepare("UPDATE email AS e, (SELECT id FROM status where name = ?) AS s SET e.status = s.id WHERE e.id = ?")
		s.Stmts["update_email_sent_at"], _ = s.Conn.Prepare("UPDATE email SET sent_at = NOW() WHERE id = ?")
		s.Stmts["update_email_opened_at"], _ = s.Conn.Prepare("UPDATE email SET opened_at = NOW() WHERE id = ?")
		s.Stmts["get_emails_by_status"], _ = s.Conn.Prepare("SELECT *, a1.email AS a1_email, a1.name AS a1_name, a2.email AS a2_email, a2.name AS a2_name FROM email AS e INNER JOIN address AS a1 ON e.`from` = a1.id INNER JOIN address AS a2 ON e.`to` = a2.id INNER JOIN status AS s ON e.status = s.id WHERE s.name = ?")
	}

	return err
}

// SaveEmail func
func (s *MySQLClient) SaveEmail(email *shared.Email) error {
	var (
		err error
		res sql.Result
	)

	email.From.IsSender = true
	err = s.getAddressID(email.From)
	if err != nil {
		return err
	}
	err = s.getAddressID(email.To)
	if err != nil {
		return err
	}

	res, err = s.Stmts["insert_email"].Exec(email.From.ID, email.To.ID, email.Subject, email.Body, email.URLUnsubscribe)
	if err != nil {
		return err
	}

	email.ID, err = res.LastInsertId()

	return err
}

// UpdateStatus func
func (s *MySQLClient) UpdateStatus(email *shared.Email, status string) error {
	_, err := s.Stmts["update_email_status"].Exec(status, email.ID)
	return err
}

// UpdateSentAt func
func (s *MySQLClient) UpdateSentAt(email *shared.Email) error {
	_, err := s.Stmts["update_email_sent_at"].Exec(email.ID)
	return err
}

// UpdateOpenedAt func
func (s *MySQLClient) UpdateOpenedAt(email *shared.Email) error {
	_, err := s.Stmts["update_email_opened_at"].Exec(email.ID)
	return err
}

// GetEmails func
func (s *MySQLClient) GetEmails(status string) ([]shared.Email, error) {
	var (
		err            error
		rows           *sql.Rows
		emails         []shared.Email
		id             sql.NullInt64
		fromID         sql.NullInt64
		fromEmail      sql.NullString
		fromName       sql.NullString
		toID           sql.NullInt64
		toEmail        sql.NullString
		toName         sql.NullString
		subject        sql.NullString
		body           sql.NullString
		urlUnsubscribe sql.NullString
		statusID       sql.NullInt64
		createdAt      *time.Time
		sentAt         *time.Time
		openedAt       *time.Time
	)

	rows, err = s.Stmts["get_emails_by_status"].Query(status)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		scanErr := rows.Scan(&id, &fromID, &toID, &subject, &body, &urlUnsubscribe, &status, &createdAt, &sentAt, &openedAt, &fromEmail, &fromName, &toEmail, &toName)

		if scanErr == nil {
			email := shared.Email{
				ID: id.Int64,
				From: &shared.Address{
					ID:    fromID.Int64,
					Email: fromEmail.String,
					Name:  fromName.String,
				},
				To: &shared.Address{
					ID:    toID.Int64,
					Email: toEmail.String,
					Name:  toName.String,
				},
				Subject:        subject.String,
				Body:           body.String,
				URLUnsubscribe: urlUnsubscribe.String,
				Status: &shared.Status{
					ID: statusID.Int64,
				},
				CreatedAt: createdAt,
				SentAt:    sentAt,
				OpenedAt:  openedAt,
			}
			emails = append(emails, email)
		}
	}

	return emails, nil
}

func (s *MySQLClient) getAddressID(address *shared.Address) error {
	isSenderInt := 0
	if address.IsSender {
		isSenderInt = 1
	}

	var (
		addressID int64
		err       error
		res       sql.Result
	)

	row := s.Stmts["get_address_id"].QueryRow(address.Email, isSenderInt)
	if row != nil {
		row.Scan(&addressID)
	}

	if addressID == 0 {
		res, err = s.Stmts["insert_address"].Exec(address.Email, address.Name, isSenderInt)
		if err != nil {
			return err
		}

		addressID, err = res.LastInsertId()
		if err != nil {
			return err
		}
	}

	address.ID = addressID

	return nil
}

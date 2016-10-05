package db

import (
	"database/sql"
	// Commented for lint
	_ "github.com/go-sql-driver/mysql"
	"github.com/logpacker/mailer/pkg/shared"
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
		s.Stmts["insert_email"], _ = s.Conn.Prepare("INSERT INTO email (`from`, `to`, subject, html) VALUES (?, ?, ?, ?)")
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

	res, err = s.Stmts["insert_email"].Exec(email.From.ID, email.To.ID, email.Subject, email.HTML)
	if err != nil {
		return err
	}

	email.ID, err = res.LastInsertId()

	return err
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

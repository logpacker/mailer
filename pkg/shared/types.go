package shared

import "time"

// MailerConfig struct
type MailerConfig struct {
	SMTPAddr       string
	MySQLAddr      string
	BeanstalkdAddr string
	APIPublicProxy string
}

// Email struct
type Email struct {
	ID             int64      `json:"id,omitempty"`
	From           *Address   `json:"from"`
	To             *Address   `json:"to"`
	Subject        string     `json:"subject"`
	Body           string     `json:"body"`
	URLUnsubscribe string     `json:"url_unsubscribe"`
	Status         *Status    `json:"status,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	OpenedAt       *time.Time `json:"opened_at,omitempty"`
}

// OpenEmail struct
type OpenEmail struct {
	ID int64 `json:"id,omitempty"`
}

// Address struct
type Address struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IsSender bool   `json:"is_sender,omitempty"`
}

// Status struct
type Status struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

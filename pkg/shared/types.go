package shared

import "time"

// Email struct
type Email struct {
	ID             int64      `json:"id,omitempty"`
	From           *Address   `json:"from"`
	To             *Address   `json:"to"`
	Subject        string     `json:"subject"`
	Body           string     `json:"body"`
	URLUnsubscribe string     `json:"url_unsubscribe"`
	Status         int64      `json:"status,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	OpenedAt       *time.Time `json:"opened_at,omitempty"`
}

// Address struct
type Address struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IsSender bool   `json:"is_sender,omitempty"`
}

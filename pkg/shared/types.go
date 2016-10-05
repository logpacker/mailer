package shared

// Email struct
type Email struct {
	ID      int64    `json:"id,omitempty"`
	From    *Address `json:"from"`
	To      *Address `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// Address struct
type Address struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IsSender bool   `json:"is_sender,omitempty"`
}

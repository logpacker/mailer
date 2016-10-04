package api

type accessParams struct {
	// API secure key
	// in: query
	// Required: true
	APIKey string `json:"api_key"`
	// Access Token
	// in: query
	// Required: true
	Token string `json:"token"`
}

// swagger:parameters tokenParams
type tokenParams struct {
	// API secure key
	// in: query
	// Required: true
	APIKey string `json:"api_key"`
}

// swagger:parameters sendParams
type sendParams struct {
	// API secure key
	// in: query
	// Required: true
	APIKey string `json:"api_key"`
	// Access Token
	// in: query
	// Required: true
	Token string `json:"token"`
	// Body parameters
	// in: body
	// Required: true
	Body sendBody `json:"body"`
}

type sendBody struct {
	From    *address `json:"from"`
	To      *address `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type address struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// swagger:response errorResponse
type errorResponse struct {
	Message string `json:"message"`
}

// swagger:response tokenResponse
type tokenResponse struct {
	Token string `json:"token"`
}

// swagger:response sendResponse
type sendResponse struct {
	ID string `json:"id"`
}

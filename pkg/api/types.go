package api

// swagger:parameters tokenParams
type tokenParams struct {
	// API secure key
	// in: query
	// Required: true
	APIKey string `json:"api_key"`
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
type sendResponse struct{}

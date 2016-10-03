package api

// swagger:parameters tokenParams
type tokenParams struct {
	// API secure key
	// in: query
	// Required: true
	APIKey string `json:"api_key"`
}
type errorResponse struct {
	Message string `json:"message"`
}

type accessTokenResponse struct {
	Token string `json:"token"`
}

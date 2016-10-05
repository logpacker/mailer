package api

import (
	"github.com/logpacker/mailer/pkg/shared"
)

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
	Body shared.Email `json:"body"`
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
	ID int64 `json:"id,omitempty"`
}

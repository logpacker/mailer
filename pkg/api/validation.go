package api

import (
	"fmt"
	"net/url"
	"regexp"
)

var regexpEmail = regexp.MustCompile(".+@.+\\..+")

func validateAccessParams(p accessParams) error {
	// Skip validation if API started without token verification
	if validAPIKey == "" {
		return nil
	}

	if p.APIKey == "" {
		return fmt.Errorf("'api_key' param is empty. API requests are forbidden")
	}
	if p.Token == "" {
		return fmt.Errorf("'token' param is empty. Get new token via 'GET /v1/token'")
	}

	tokensMu.Lock()
	t, _ := tokens[p.APIKey]
	tokensMu.Unlock()

	if t != p.Token {
		return fmt.Errorf("'token' is not valid. Get new token via 'GET /v1/token'")
	}

	return nil
}

func validateTokenParams(p tokenParams) error {
	// Skip validation if API started without token verification
	if validAPIKey == "" {
		return nil
	}

	if p.APIKey != validAPIKey {
		return fmt.Errorf("'api_key' is not valid")
	}

	return nil
}

func validateSendParams(p sendParams) error {
	if p.Body.From == nil {
		return fmt.Errorf("'body.from' is empty")
	}
	if p.Body.To == nil {
		return fmt.Errorf("'body.to' is empty")
	}
	if p.Body.Subject == "" {
		return fmt.Errorf("'body.subject' is empty")
	}
	if p.Body.Body == "" {
		return fmt.Errorf("'body.body' is empty")
	}
	if p.Body.URLUnsubscribe == "" {
		return fmt.Errorf("'body.url_unsubscribe' is empty")
	}

	if p.Body.From.Email == "" {
		return fmt.Errorf("'body.from.email' is empty")
	}
	if p.Body.From.Name == "" {
		return fmt.Errorf("'body.from.name' is empty")
	}
	if p.Body.To.Email == "" {
		return fmt.Errorf("'body.to.email' is empty")
	}

	if !regexpEmail.Match([]byte(p.Body.From.Email)) {
		return fmt.Errorf("'body.from.email' is invalid")
	}
	if !regexpEmail.Match([]byte(p.Body.To.Email)) {
		return fmt.Errorf("'body.to.email' is invalid")
	}

	urlRes, urlErr := url.Parse(p.Body.URLUnsubscribe)
	if urlErr != nil || urlRes.Host == "" {
		return fmt.Errorf("'body.url_unsubscribe' is not a valid URL")
	}

	return nil
}

func validateTrackParams(p trackParams) error {
	if p.ID < 1 {
		return fmt.Errorf("'id' cannot be empty")
	}

	return nil
}

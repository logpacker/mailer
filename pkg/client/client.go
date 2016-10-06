package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Config struct. URL and APIKey are mandatory for initialization
type Config struct {
	URL    string
	APIKey string
	Token  *string
}

// Client struct
type Client struct {
	Config Config
}

// Address struct
type Address struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Email struct. Body is a HTML inside <body></body>
type Email struct {
	From           *Address `json:"from"`
	To             *Address `json:"to"`
	Subject        string   `json:"subject"`
	Body           string   `json:"body"`
	URLUnsubscribe string   `json:"url_unsubscribe"`
}

// New client constructor
func New(conf Config) (*Client, error) {
	if conf.URL == "" {
		return nil, fmt.Errorf("Config.URL cannot be empty")
	}

	c := new(Client)
	c.Config = conf

	return c, nil
}

type errorResponse struct {
	Message string `json:"message"`
}
type tokenResponse struct {
	Token string `json:"token"`
}
type sendResponse struct {
	Status bool `json:"status"`
}

// GetToken returns cached token or make a new request to get it by APIKey
func (c *Client) GetToken() (string, error) {
	if c == nil {
		return "", fmt.Errorf("Client is nil")
	}

	if c.Config.Token != nil {
		return *c.Config.Token, nil
	}

	tr := tokenResponse{}
	err := c.api("/v1/token", "GET", nil, &tr)
	if err != nil {
		return "", err
	}

	c.Config.Token = &tr.Token
	if c.Config.Token != nil {
		return *c.Config.Token, nil
	}

	return "", fmt.Errorf("Token is nil")
}

// SendEmail sends email asynchronously
func (c *Client) SendEmail(email Email) error {
	if c == nil {
		return fmt.Errorf("Client is nil")
	}

	var (
		err  error
		body []byte
	)

	if c.Config.Token == nil {
		_, err = c.GetToken()
		if err != nil {
			return err
		}
	}

	body, err = json.Marshal(email)
	if err != nil {
		return err
	}

	sr := sendResponse{}
	err = c.api("/v1/send", "POST", body, &sr)
	return err
}

func (c *Client) api(endpoint string, method string, bodyData []byte, v interface{}) error {
	if c == nil {
		return fmt.Errorf("Client is nil")
	}

	var (
		err      error
		body     []byte
		request  *http.Request
		response *http.Response
	)

	// Build request
	t := ""
	if c.Config.Token != nil {
		t = *c.Config.Token
	}
	u := fmt.Sprintf("%s/%s?api_key=%s&token=%s", strings.TrimRight(c.Config.URL, "/"), strings.TrimLeft(endpoint, "/"), c.Config.APIKey, t)
	request, err = http.NewRequest(method, u, bytes.NewBuffer(bodyData))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", fmt.Sprintf("%d", len(bodyData)))
	if err != nil {
		return err
	}

	// Do request and get raw body
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err = client.Do(request)
	if err != nil {
		return err
	}

	body, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}

	errResp := errorResponse{}
	err = json.Unmarshal(body, &errResp)
	if err == nil && errResp.Message != "" {
		return fmt.Errorf("Mailer API error. URL: %s. Details: %s", u, errResp.Message)
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("Unable to parse response from Mailer API. URL: %s. Details: %s", u, err.Error())
		}
	}

	return nil
}

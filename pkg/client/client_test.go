package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server    *httptest.Server
	testToken = "testtoken"
)

func init() {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var js []byte
		if r.URL.Path == "/v1/token" {
			js, _ = json.Marshal(tokenResponse{
				Token: testToken,
			})
		}
		if r.URL.Path == "/v1/send" {
			js, _ = json.Marshal(sendResponse{
				Status: true,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}))
}

func TestConfig(t *testing.T) {
	token := testToken

	c := Config{
		server.URL,
		"key",
		&token,
	}

	if c.URL != server.URL {
		t.Errorf("Config struct error")
	}
	if c.APIKey != "key" {
		t.Errorf("Config struct error")
	}
	if c.Token != &token {
		t.Errorf("Config struct error")
	}
}

func TestClient(t *testing.T) {
	c := Client{
		Config{
			URL: server.URL,
		},
	}

	if c.Config.URL != server.URL {
		t.Errorf("Client struct error")
	}
}

func TestAddress(t *testing.T) {
	a := Address{
		"foo@bar.com",
		"name",
	}

	if a.Email != "foo@bar.com" {
		t.Errorf("Address struct error")
	}
	if a.Name != "name" {
		t.Errorf("Address struct error")
	}
}

func TestEmail(t *testing.T) {
	e := Email{
		&Address{
			Email: "from@bar.com",
		},
		&Address{
			Email: "to@bar.com",
		},
		"subject",
		"body",
		"unsubscribe_url",
	}

	if e.From.Email != "from@bar.com" {
		t.Errorf("Email struct error")
	}
	if e.To.Email != "to@bar.com" {
		t.Errorf("Email struct error")
	}
	if e.Subject != "subject" {
		t.Errorf("Email struct error")
	}
	if e.Body != "body" {
		t.Errorf("Email struct error")
	}
	if e.URLUnsubscribe != "unsubscribe_url" {
		t.Errorf("Email struct error")
	}
}

func TestNew(t *testing.T) {
	conf1 := Config{}
	conf2 := Config{
		URL: server.URL,
	}

	c1, err1 := New(conf1)
	if err1 == nil || c1 != nil {
		t.Errorf("New func error")
	}

	c2, err2 := New(conf2)
	if err2 != nil || c2 == nil {
		t.Errorf("New func error")
	}
}

func TestGetToken(t *testing.T) {
	c, _ := New(Config{
		URL: server.URL,
	})

	t1, e1 := c.GetToken()
	if e1 != nil || t1 != testToken || *c.Config.Token != testToken {
		t.Errorf("GetToken func error")
	}

	t2, e2 := c.GetToken()
	if e2 != nil || t2 != testToken || *c.Config.Token != testToken {
		t.Errorf("GetToken func error")
	}
}

func TestSendEmail(t *testing.T) {
	c, _ := New(Config{
		URL: server.URL,
	})

	err := c.SendEmail(Email{})
	if err != nil {
		t.Errorf("SendEmail func error")
	}
}

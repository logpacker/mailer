package api

import (
	"encoding/json"
	"github.com/gocraft/web"
	"net/http"
	"sync"
)

var (
	tokensMu    sync.Mutex
	tokens      map[string]string
	validAPIKey string
)

// Context type
type Context struct {
	APIKey string
	Token  string
}

func (c *Context) checkToken(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.APIKey = r.FormValue("api_key")
	c.Token = r.FormValue("token")
	if c.APIKey == "" {
		c.writeResponse(w, r, errorResponse{
			Message: "'api_key' param is empty. API requests are forbidden",
		})
		return
	}

	if r.URL.Path != "/v1/token" {
		if c.Token == "" {
			c.writeResponse(w, r, errorResponse{
				Message: "'token' param is empty. Get new token via 'GET /v1/token'",
			})
			return
		}

		tokensMu.Lock()
		t, _ := tokens[c.APIKey]
		tokensMu.Unlock()

		if t != c.Token {
			c.writeResponse(w, r, errorResponse{
				Message: "'token' is not valid. Get new token via 'GET /v1/token'",
			})
			return
		}
	}

	next(w, r)
}

func (c *Context) writeResponse(w web.ResponseWriter, r *web.Request, response interface{}) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// NewRouter - constructor
func NewRouter(apiKey string) *web.Router {
	validAPIKey = apiKey

	tokens = make(map[string]string)
	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware((*Context).checkToken).
		Get("/v1/token", (*Context).token)

	return router
}

// swagger:route GET /v1/token tokenParams
//
// Get access token
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       401: errorResponse
//       200: accessTokenResponse
func (c *Context) token(w web.ResponseWriter, r *web.Request) {
	if c.APIKey != validAPIKey {
		c.writeResponse(w, r, errorResponse{
			Message: "'api_key' is not valid",
		})
		return
	}

	tokensMu.Lock()
	tokens[c.APIKey] = "token"
	tokensMu.Unlock()

	c.writeResponse(w, r, accessTokenResponse{
		Token: tokens[c.APIKey],
	})
}

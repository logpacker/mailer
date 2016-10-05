package api

import (
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/logpacker/mailer/pkg/conf"
	"github.com/logpacker/mailer/pkg/db"
	"github.com/logpacker/mailer/pkg/queue"
	"github.com/logpacker/mailer/pkg/shared"
	"net/http"
	"sync"
)

var (
	tokensMu    sync.Mutex
	tokens      map[string]string
	validAPIKey string
	dbClient    *db.MySQLClient
	queueClient *queue.BeanstalkdClient
)

// Context type
type Context struct{}

func (c *Context) checkToken(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	r.ParseForm()

	if r.URL.Path != "/v1/token" {
		p := accessParams{
			APIKey: r.FormValue("api_key"),
			Token:  r.FormValue("token"),
		}

		if err := validateAccessParams(p); err != nil {
			c.writeErrorResponse(w, r, err)
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

func (c *Context) writeErrorResponse(w web.ResponseWriter, r *web.Request, err error) {
	if err != nil {
		c.writeResponse(w, r, errorResponse{
			Message: err.Error(),
		})
	}
}

// NewRouter - constructor
func NewRouter(apiKey string, conf *conf.MailerConfig) *web.Router {
	validAPIKey = apiKey
	tokens = make(map[string]string)
	dbClient = new(db.MySQLClient)
	dbErr := dbClient.Init(conf.MySQLAddr)
	if dbErr != nil {
		panic(dbErr)
	}

	queueClient = new(queue.BeanstalkdClient)
	queueErr := queueClient.Init(conf.BeanstalkdAddr)
	if queueErr != nil {
		panic(queueErr)
	}

	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware((*Context).checkToken).
		Get("/v1/token", (*Context).token).
		Post("/v1/send", (*Context).send)

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
//       200: tokenResponse
func (c *Context) token(w web.ResponseWriter, r *web.Request) {
	p := tokenParams{
		APIKey: r.FormValue("api_key"),
	}

	if err := validateTokenParams(p); err != nil {
		c.writeErrorResponse(w, r, err)
		return
	}

	tokensMu.Lock()
	tokens[p.APIKey] = "token"
	tokensMu.Unlock()

	c.writeResponse(w, r, tokenResponse{
		Token: tokens[p.APIKey],
	})
}

// swagger:route POST /v1/send sendParams
//
// Adds email into the mailer queue
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
//       200: sendResponse
func (c *Context) send(w web.ResponseWriter, r *web.Request) {
	decoder := json.NewDecoder(r.Body)
	var b shared.Email
	jsonErr := decoder.Decode(&b)
	if jsonErr != nil {
		c.writeErrorResponse(w, r, jsonErr)
		return
	}

	p := sendParams{
		Body: b,
	}

	if err := validateSendParams(p); err != nil {
		c.writeErrorResponse(w, r, err)
		return
	}

	prepareErr := prepareEmail(&b)
	if prepareErr != nil {
		c.writeErrorResponse(w, r, prepareErr)
		return
	}

	saveErr := dbClient.SaveEmail(&b)
	if saveErr != nil {
		c.writeErrorResponse(w, r, saveErr)
		return
	}

	queueErr := queueClient.SendEmailJob(&b)
	if queueErr != nil {
		c.writeErrorResponse(w, r, queueErr)
		return
	}

	c.writeResponse(w, r, sendResponse{
		ID: b.ID,
	})
}

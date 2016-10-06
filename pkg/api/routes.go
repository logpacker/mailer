package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/logpacker/mailer/pkg/queue"
	"github.com/logpacker/mailer/pkg/shared"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	tokensMu    sync.Mutex
	tokens      map[string]string
	validAPIKey string
	queueClient *queue.BeanstalkdClient
)

// Context type
type Context struct{}

func (c *Context) checkToken(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	r.ParseForm()

	if r.URL.Path != "/v1/token" && r.URL.Path != "/v1/track" {
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
func NewRouter(apiKey string, conf *shared.MailerConfig) *web.Router {
	validAPIKey = apiKey
	tokens = make(map[string]string)

	queueClient = new(queue.BeanstalkdClient)
	queueErr := queueClient.Init(conf.BeanstalkdAddr)
	if queueErr != nil {
		panic(queueErr)
	}

	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware((*Context).checkToken).
		Get("/v1/token", (*Context).token).
		Post("/v1/send", (*Context).send).
		Get("/v1/track", (*Context).track)

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
	_, exists := tokens[p.APIKey]
	if !exists {
		h2 := sha1.New()
		h2.Write([]byte(p.APIKey + time.Now().Format("20060102150405")))
		tokens[p.APIKey] = hex.EncodeToString(h2.Sum(nil))
	}
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

	queueErr := queueClient.SendEmailJob(&b)
	if queueErr != nil {
		c.writeErrorResponse(w, r, queueErr)
		return
	}

	c.writeResponse(w, r, sendResponse{
		Status: true,
	})
}

// swagger:route POST /v1/track trackParams
//
// Mark email as opened by Client when tracker image is loaded
//
//     Produces:
//     - image/png
//
//     Schemes: http
func (c *Context) track(w web.ResponseWriter, r *web.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	p := trackParams{
		ID: int64(id),
	}

	validateErr := validateTrackParams(p)
	shared.LogErr(validateErr)

	queueErr := queueClient.SendOpenJob(&shared.OpenEmail{
		ID: p.ID,
	})
	shared.LogErr(queueErr)

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(1, 1, color.Transparent)

	png.Encode(w, img)
	w.Header().Set("Content-Type", "image/png")
}

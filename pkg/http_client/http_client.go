package httpclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	AcceptJson       = "application/json"
	AcceptRest       = "application/vnd.pgrst.object+json"
	ServiceUserAgent = "service_template"
	ContentType      = "application/json; charset=utf-8"
)

type HttpClient struct {
	BaseURI       string
	TransactionID string
	UserAgent     string
	ContentType   string
	PrivateToken  string
	Accept        string
	Timeout       time.Duration
	Debug         bool
}

func Default() *HttpClient {
	return &HttpClient{
		TransactionID: time.Now().String(),
		UserAgent:     ServiceUserAgent,
		ContentType:   ContentType,
		Accept:        AcceptJson,
		Debug:         true,
	}
}

func (c *HttpClient) FastGet(requestURI string) (*fasthttp.Response, error) {
	t1 := time.Now()
	c.TransactionID = "123"

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()
	request.SetRequestURI(fmt.Sprintf("%s%s", c.BaseURI, requestURI))
	request.Header.SetContentType(c.ContentType)
	request.Header.Add(fasthttp.HeaderUserAgent, c.UserAgent)
	request.Header.Add("Transaction-Id", c.TransactionID)
	request.Header.Add(fasthttp.HeaderAccept, c.Accept)
	request.Header.Add("apikey", c.PrivateToken)
	request.Header.SetMethod(fasthttp.MethodGet)

	if c.Debug {
		request.Header.VisitAll(func(key, value []byte) {
			logrus.WithFields(logrus.Fields{
				"Transaction-ID": c.TransactionID,
				string(key):      string(value),
			}).Debug("Http-Client-Request")
		})
	}

	timeout := time.Second * 3
	if c.Timeout != 0 {
		timeout = c.Timeout
	}

	err := fasthttp.DoTimeout(request, response, timeout)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		elapsed := time.Since(t1)
		response.Header.VisitAll(func(key, value []byte) {
			logrus.WithFields(logrus.Fields{
				"Transaction-ID": c.TransactionID,
				"elapsed":        elapsed.String(),
				"Status-Code":    strconv.Itoa(response.StatusCode()),
				string(key):      string(value),
				"payload":        string(response.Body()),
			}).Debug("Http-Client-Response")
		})

	}

	out := fasthttp.AcquireResponse()
	response.CopyTo(out)

	return out, nil
}

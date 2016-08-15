package caddydogstatsd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/datadog/datadog-go/statsd"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// DogstatsdHandler is a middleware handler for reporting dogstatsd metrics on requests
type DogstatsdHandler struct {
	Client     *statsd.Client
	SampleRate float64
	Next       httpserver.Handler
}

// ServeHTTP is the middleware handler which will emit dogstatsd metrics after handling a request
func (h DogstatsdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// If we do not have a statsd.Client configured, then skip any processing
	if h.Client == nil {
		return h.Next.ServeHTTP(w, r)
	}

	// Grab the request start time
	var start time.Time
	start = time.Now()

	// Handle the request
	var code int
	var err error
	code, err = h.Next.ServeHTTP(w, r)

	// Grab the request durection in Milliseconds
	var elapsed time.Duration
	var elapsedMS float64
	elapsed = time.Since(start)
	elapsedMS = float64(elapsed.Nanoseconds()) / float64(time.Millisecond)

	// Report our request metrics to dogstatsd
	var client statsd.Client
	client = *h.Client
	var extraTags = []string{
		fmt.Sprintf("status_code:%d", code),
	}
	client.Count("caddy.response.count", 1, extraTags, h.SampleRate)
	client.TimeInMilliseconds("caddy.response.time", elapsedMS, extraTags, h.SampleRate)
	return code, err
}

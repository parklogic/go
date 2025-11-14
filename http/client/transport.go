package client

import (
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type transport struct {
	transport http.Transport
	userAgent string
}

func newTransport(cfg *Configuration) *transport {
	return &transport{
		transport: http.Transport{
			DialContext: (&net.Dialer{
				FallbackDelay: -1,
				KeepAlive:     cfg.KeepAlive,
				Timeout:       cfg.ConnectionTimeout,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			IdleConnTimeout:       cfg.IdleConnTimeout,
			MaxIdleConnsPerHost:   cfg.MaxIdleConnPerHost,
			ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
			TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
		},
		userAgent: cfg.UserAgent,
	}
}

func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	logger := zerolog.Ctx(req.Context()).With().Str("method", req.Method).Str("url", req.URL.String()).Logger()

	var start time.Time
	if logger.GetLevel() <= zerolog.TraceLevel {
		start = time.Now()
	}
	logger.Trace().Msg("Sending HTTP request")
	defer func() {
		var took time.Duration
		if logger.GetLevel() <= zerolog.TraceLevel {
			took = time.Since(start)
		}

		if err == nil {
			logger.Trace().Dur("took", took).Str("status", res.Status).Msg("Got HTTP response")
		}
	}()

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", t.userAgent)
	}

	return t.transport.RoundTrip(req)
}

package api

import (
	"fmt"
	"net/http"
)

type HandlerFunc[P any, R any] func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error)

func HandlerMethodNotAllowed(_ http.ResponseWriter, r *http.Request, _ any) (status int, res *any, err error) {
	return 0, nil, NewProblemDetailFromStatus(http.StatusMethodNotAllowed, ErrMethodNotAllowed, fmt.Sprintf("Request method %q not allowed", r.Method))
}

func HandlerNotFound(_ http.ResponseWriter, _ *http.Request, _ any) (status int, res *any, err error) {
	return 0, nil, NewProblemDetailFromStatus(http.StatusNotFound, ErrEndpointNotFound, "Invalid API endpoint")
}

// HandlerLive responds with a JSON containing a static status message.
//
//	@Summary	Get liveness status
//	@Success	200	{object}	Status
//	@Tags		server
//	@Router		/livez [get]
func HandlerLive(_ http.ResponseWriter, _ *http.Request, _ NilPayload) (status int, res *Status, err error) {
	return http.StatusOK, &Status{Status: "ok"}, nil
}

type NilPayload struct{}

// Status represents a server status response
type Status struct {
	Status string `json:"status" example:"ok"`
} //	@Name	ServerStatus

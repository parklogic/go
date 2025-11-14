package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	errPkg "github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func Adapter[P any, R any](h HandlerFunc[P, R]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, body, err := execHandlerFunc[P, R](w, r, h)
		if err != nil {
			handlerErr(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if body == nil || string(body) == "null" {
			return
		}

		if _, err = w.Write(body); err != nil {
			logger := zerolog.Ctx(r.Context())

			logger.Error().Err(err).Msg("Error writing response")

			return
		}
	}
}

func decodePayload[T any](r *http.Request, v *T) error {
	if r.ContentLength == 0 {
		return nil
	}

	switch any(v).(type) {
	case nil:
		return nil
	case *NilPayload:
		return ErrPayloadNotExpected
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(v)
}

func execHandlerFunc[P any, R any](w http.ResponseWriter, r *http.Request, h HandlerFunc[P, R]) (status int, body []byte, err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			if rvr == http.ErrAbortHandler {
				panic(rvr)
			}

			if e, ok := rvr.(error); ok {
				err = e
				return
			}

			if s, ok := rvr.(string); ok {
				err = errors.New(s)
				return
			}
		}
	}()

	var p P
	var res *R

	err = decodePayload(r, &p)
	if err != nil {
		return 0, nil, err
	}

	status, res, err = h(w, r, p)
	if err != nil {
		return 0, nil, err
	}

	if res == nil {
		return status, nil, nil
	}

	body, err = json.Marshal(res)
	if err != nil {
		return 0, nil, err
	}

	return status, body, nil
}

func handlerErr(w http.ResponseWriter, r *http.Request, err error) {
	logger := zerolog.Ctx(r.Context()).With().Err(errPkg.WithStack(err)).Logger()

	// todo: fix this
	// logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
	// 	return c.Err(errPkg.WithStack(err))
	// })

	if errors.Is(err, context.Canceled) {
		return
	}

	var apiErr Error
	if !errors.As(err, &apiErr) {
		apiErr = NewErrorFromErr(err)
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(apiErr.StatusCode())

	body, err := json.Marshal(apiErr)
	if err != nil {
		logger.Error().Err(err).Msg("Error marshalling error payload")
		return
	}

	if _, err = w.Write(body); err != nil {
		logger.Error().Err(err).Msg("Error writing error response")
		return
	}
}

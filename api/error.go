package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/parklogic/go/bearer"
	"github.com/parklogic/go/database"
	"github.com/parklogic/go/pagination"
	"github.com/parklogic/go/value"
)

var (
	ErrEndpointNotFound   = fmt.Errorf("endpoint not found")
	ErrMethodNotAllowed   = fmt.Errorf("request method not allowed")
	ErrPayloadNotExpected = fmt.Errorf("no payload expected")
)

type Error interface {
	Error() string

	StatusCode() int
}

func NewErrorFromErr(err error) Error {
	switch {
	case errors.Is(err, ErrPayloadNotExpected):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, "No payload expected")

	case errors.Is(err, bearer.ErrInvalidAuthType):
		return NewProblemDetailFromStatus(http.StatusUnauthorized, err, "Invalid Authorization type")
	case errors.Is(err, bearer.ErrInvalidToken):
		return NewProblemDetailFromStatus(http.StatusForbidden, err, "Invalid bearer token")
	case errors.Is(err, bearer.ErrMissingHeader):
		return NewProblemDetailFromStatus(http.StatusUnauthorized, err, "Missing authorization header")
	case errors.Is(err, bearer.ErrMissingToken):
		return NewProblemDetailFromStatus(http.StatusUnauthorized, err, "Missing bearer token")

	case errors.Is(err, database.ErrAlreadyExists):
		return NewProblemDetailFromStatus(http.StatusConflict, err, "Resource unique value already exists")
	case errors.Is(err, database.ErrForeignKeyConstrainError):
		return NewProblemDetailFromStatus(http.StatusConflict, err, "Failed to validate related resource constrain")
	case errors.Is(err, database.ErrNotFound):
		return NewProblemDetailFromStatus(http.StatusNotFound, err, "Resource not found")

	case errors.Is(err, pagination.ErrInvalidLimit):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, "Invalid pagination limit")
	case errors.Is(err, pagination.ErrInvalidPage):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, "Invalid pagination page number")
	case errors.Is(err, pagination.ErrInvalidSortKey):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, "Invalid pagination sort key")
	case errors.Is(err, pagination.ErrInvalidSortOrder):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, "Invalid pagination sort order")
	case errors.Is(err, pagination.ErrMissingSortKey):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, "Missing pagination sort key")
	}

	if errVal := (validator.ValidationErrors{}); errors.As(err, &errVal) {
		apiErr := ValidationProblemDetail{
			ProblemDetail: NewProblemDetailFromStatus(http.StatusBadRequest, err, "Validation error"),
			InvalidParams: make([]ValidationInvalidParam, len(errVal)),
		}

		for i, fe := range errVal {
			apiErr.InvalidParams[i] = ValidationInvalidParam{
				Name:   fe.Field(),
				Reason: "Field validation failed",
				Rule: ValidationRule{
					Name:      fe.Tag(),
					Parameter: fe.Param(),
				},
			}
		}

		return apiErr
	}

	if errVal := (&value.Error{}); errors.As(err, &errVal) {
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, err.Error())
	}

	if errVal := (&strconv.NumError{}); errors.As(err, &errVal) {
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, fmt.Sprintf("Invalid value %q: %s", errVal.Num, errVal.Err))
	}

	if errVal := (&json.SyntaxError{}); errors.As(err, &errVal) {
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, fmt.Sprintf("Invalid JSON syntax: %s", errVal))
	}

	if errVal := (&json.UnmarshalTypeError{}); errors.As(err, &errVal) {
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, fmt.Sprintf("Received value of type %q, should be %q", errVal.Value, errVal.Type))
	}

	errMsg := err.Error()
	switch {
	case strings.HasPrefix(errMsg, "json: unknown field "):
		return NewProblemDetailFromStatus(http.StatusBadRequest, err, fmt.Sprintf("Unknown field %s", errMsg[20:]))
	}

	return NewProblemDetailFromStatus(http.StatusInternalServerError, err, "Internal server error")
}

// ProblemDetail represents an API problem detail.
type ProblemDetail struct {
	// Specific explanation of the error
	Detail string `json:"detail,omitempty" example:"The requested resource was not found"`
	// URI identifier for a specific occurrence of the problem
	Instance string `json:"instance,omitempty" example:"/user/1"`
	// HTTP status code
	Status int `json:"status" example:"404"`
	// Summary of the error
	Title string `json:"title" example:"Not Found"`
	// URI identifier for the problem type
	Type string `json:"type,omitempty" example:"https://example.net/resource-not-found"`

	err error
} //	@name	ProblemDetail

func NewProblemDetailFromStatus(status int, err error, detail string) ProblemDetail {
	errTitle := http.StatusText(status)
	errType := fmt.Sprintf("/docs/error/%s", url.PathEscape(strings.ToLower(errTitle)))

	return ProblemDetail{
		Detail: detail,
		Status: status,
		Title:  errTitle,
		Type:   errType,

		err: err,
	}
}

func (e ProblemDetail) Error() string {
	return e.err.Error()
}

func (e ProblemDetail) StatusCode() int {
	return e.Status
}

func (e ProblemDetail) Unwrap() error {
	return e.err
}

// ValidationProblemDetail represents an API validation problem detail.
type ValidationProblemDetail struct {
	ProblemDetail

	InvalidParams []ValidationInvalidParam `json:"invalid_params"`
}

func (e ValidationProblemDetail) Unwrap() error {
	return e.ProblemDetail
}

type ValidationInvalidParam struct {
	Name   string         `json:"name"`
	Reason string         `json:"reason"`
	Rule   ValidationRule `json:"rule"`
}

type ValidationRule struct {
	Name      string `json:"name,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

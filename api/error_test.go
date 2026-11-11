package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/parklogic/go/bearer"
	"github.com/parklogic/go/database"
	"github.com/parklogic/go/pagination"
	"github.com/parklogic/go/value"
	"github.com/stretchr/testify/require"
)

type stubAPIError struct{}

func (stubAPIError) Error() string { return "boom" }

type requiredName struct {
	Name string `validate:"required"`
}

func TestNewErrorFromErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "payload", err: ErrPayloadNotExpected, status: http.StatusBadRequest},
		{name: "auth type", err: bearer.ErrInvalidAuthType, status: http.StatusUnauthorized},
		{name: "token", err: bearer.ErrInvalidToken, status: http.StatusForbidden},
		{name: "header", err: bearer.ErrMissingHeader, status: http.StatusUnauthorized},
		{name: "missing token", err: bearer.ErrMissingToken, status: http.StatusUnauthorized},
		{name: "exists", err: database.ErrAlreadyExists, status: http.StatusConflict},
		{name: "foreign key", err: database.ErrForeignKeyConstrainError, status: http.StatusConflict},
		{name: "missing resource", err: database.ErrNotFound, status: http.StatusNotFound},
		{name: "limit", err: pagination.ErrInvalidLimit, status: http.StatusBadRequest},
		{name: "page", err: pagination.ErrInvalidPage, status: http.StatusBadRequest},
		{name: "sort key", err: pagination.ErrInvalidSortKey, status: http.StatusBadRequest},
		{name: "sort order", err: pagination.ErrInvalidSortOrder, status: http.StatusBadRequest},
		{name: "sort missing", err: pagination.ErrMissingSortKey, status: http.StatusBadRequest},
		{name: "unknown", err: stubAPIError{}, status: http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewErrorFromErr(test.err)
			require.Equal(t, test.status, got.StatusCode())
			require.ErrorIs(t, got, test.err)
		})
	}
}

func TestNewErrorFromErrValues(t *testing.T) {
	t.Parallel()

	validationErr := validator.New().Struct(requiredName{})
	require.Error(t, validationErr)

	numberErr := strconv.ErrSyntax
	_, parseErr := strconv.Atoi("nope")
	require.ErrorIs(t, parseErr, numberErr)

	var destination struct {
		Count int `json:"count"`
	}

	syntaxErr := json.Unmarshal([]byte("{"), &destination)
	typeErr := json.Unmarshal([]byte(`{"count":"x"}`), &destination)

	tests := []struct {
		name string
		err  error
	}{
		{name: "validation", err: validationErr},
		{name: "value", err: value.NewError(stubAPIError{})},
		{name: "number", err: parseErr},
		{name: "syntax", err: syntaxErr},
		{name: "type", err: typeErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewErrorFromErr(test.err)
			require.Equal(t, http.StatusBadRequest, got.StatusCode())

			if test.name == "validation" {
				require.EqualError(t, got, test.err.Error())

				var detail ValidationProblemDetail

				require.ErrorAs(t, got, &detail)
				require.NotEmpty(t, detail.InvalidParams)

				return
			}

			require.ErrorIs(t, got, test.err)
		})
	}
}

func TestNewProblemDetailFromStatus(t *testing.T) {
	t.Parallel()

	got := NewProblemDetailFromStatus(http.StatusNotFound, stubAPIError{}, "missing")
	require.Equal(t, http.StatusNotFound, got.StatusCode())
	require.EqualError(t, got, "boom")
	require.ErrorIs(t, got, stubAPIError{})
	require.Equal(t, "Not Found", got.Title)
	require.NotEmpty(t, got.Type)
	require.Equal(t, stubAPIError{}, got.Unwrap())
}

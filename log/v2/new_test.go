package log

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Configuration
		handler string
		err     error
	}{
		{
			name:    "Format logfmt",
			cfg:     &Configuration{Format: "logfmt", Level: "info"},
			handler: "*slog.TextHandler",
		},
		{
			name:    "Format json",
			cfg:     &Configuration{Format: "json", Level: "info"},
			handler: "*slog.JSONHandler",
		},
		{
			name:    "Level debug",
			cfg:     &Configuration{Format: "logfmt", Level: "debug"},
			handler: "*slog.TextHandler",
		},
		{
			name:    "Level info",
			cfg:     &Configuration{Format: "logfmt", Level: "info"},
			handler: "*slog.TextHandler",
		},
		{
			name:    "Level warn",
			cfg:     &Configuration{Format: "logfmt", Level: "warn"},
			handler: "*slog.TextHandler",
		},
		{
			name:    "Level error",
			cfg:     &Configuration{Format: "logfmt", Level: "error"},
			handler: "*slog.TextHandler",
		},
		{
			name: "invalid format",
			cfg:  &Configuration{Format: "invalid", Level: "info"},
			err:  ErrInvalidConfig{Field: "format", Value: "invalid"},
		},
		{
			name: "invalid level",
			cfg:  &Configuration{Format: "logfmt", Level: "invalid"},
			err:  ErrInvalidConfig{Field: "level", Value: "invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := New(tt.cfg)

			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
				return

			}

			if !assert.NoError(t, err) {
				t.FailNow()
			}

			hander := reflect.TypeOf(l.Handler()).String()

			assert.Equal(t, tt.handler, hander)
		})
	}

}

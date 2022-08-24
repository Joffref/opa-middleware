package config

import (
	"github.com/open-policy-agent/opa/rego"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		URL                 string
		Policy              string
		instantiatedPolicy  *rego.PreparedEvalQuery
		Query               string
		InputCreationMethod func(r *http.Request) (map[string]interface{}, error)
		ExceptedResult      bool
		DeniedStatusCode    int
		DeniedMessage       string
		Headers             map[string]string
		IgnoredHeaders      []string
		Debug               bool
		Logger              *log.Logger
		Timeout             time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid config",
			fields: fields{
				URL:   "http://localhost:8080",
				Query: "data.test.allow",
				InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
					return map[string]interface{}{}, nil
				},
				ExceptedResult:   true,
				DeniedStatusCode: http.StatusForbidden,
				DeniedMessage:    "Forbidden",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				IgnoredHeaders: []string{
					"Content-Type",
				},
				Debug: true,
			},
			wantErr: false,
		},
		{
			name: "invalid config given URL and Policy are set",
			fields: fields{
				URL:    "http://localhost:8080",
				Policy: "policy",
				Query:  "data.test.allow",
				InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
					return map[string]interface{}{}, nil
				},
				ExceptedResult:   true,
				DeniedStatusCode: http.StatusForbidden,
				DeniedMessage:    "Forbidden",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config given URL and Policy are not set",
			fields: fields{
				Query: "data.test.allow",
				InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
					return map[string]interface{}{}, nil
				},
				ExceptedResult:   true,
				DeniedStatusCode: http.StatusForbidden,
				DeniedMessage:    "Forbidden",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config given no input creation method",
			fields: fields{
				URL:              "http://localhost:8080",
				Query:            "data.test.allow",
				ExceptedResult:   true,
				DeniedStatusCode: http.StatusForbidden,
				DeniedMessage:    "Forbidden",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				IgnoredHeaders: []string{
					"Content-Type",
				},
				Debug: true,
			},
			wantErr: true,
		},
		{
			name: "invalid config given no query",
			fields: fields{
				URL: "http://localhost:8080",
				InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
					return map[string]interface{}{}, nil
				},
				ExceptedResult:   true,
				DeniedStatusCode: http.StatusForbidden,
				DeniedMessage:    "Forbidden",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				IgnoredHeaders: []string{
					"Content-Type",
				},
				Debug: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				URL:                 tt.fields.URL,
				Policy:              tt.fields.Policy,
				instantiatedPolicy:  tt.fields.instantiatedPolicy,
				Query:               tt.fields.Query,
				InputCreationMethod: tt.fields.InputCreationMethod,
				ExceptedResult:      tt.fields.ExceptedResult,
				DeniedStatusCode:    tt.fields.DeniedStatusCode,
				DeniedMessage:       tt.fields.DeniedMessage,
				Headers:             tt.fields.Headers,
				IgnoredHeaders:      tt.fields.IgnoredHeaders,
				Debug:               tt.fields.Debug,
				Logger:              tt.fields.Logger,
				Timeout:             tt.fields.Timeout,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

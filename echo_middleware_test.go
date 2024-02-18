package opamiddleware

import (
	"github.com/Joffref/opa-middleware/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var Test_Policy = `
package policy

default allow = false

allow {
	input.path = "/api/v1/users"
	input.method = "GET"
}`

func TestEchoMiddleware_Query(t *testing.T) {
	type fields struct {
		Config *config.Config
		InputCreationMethod   EchoInputCreationMethod
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test EchoMiddleware_Query",
			fields: fields{
				Config: &config.Config{
					Policy: Test_Policy,
					Query:  "data.policy.allow",
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				InputCreationMethod: func(c echo.Context) (map[string]interface{}, error) {
					return map[string]interface{}{
						"path":   c.Request().URL.Path,
						"method": c.Request().Method,
					}, nil
				},
			},
			args: args{
				req: &http.Request{
					URL: &url.URL{
						Path: "/api/v1/users",
					},
					Method: "GET",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {		
			e := echo.New()	
			h := &EchoMiddleware{
				Config: tt.fields.Config,
				InputCreationMethod:   tt.fields.InputCreationMethod,
			}
			c := e.NewContext(tt.args.req, httptest.NewRecorder())
			got, err := h.query(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEchoMiddleware_Use(t *testing.T) {
	type fields struct {
		Config *config.Config
		InputCreationMethod   EchoInputCreationMethod
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test EchoMiddleware_Use",
			fields: fields{
				Config: &config.Config{
					Policy: Test_Policy,
					Query:  "data.policy.allow",
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				InputCreationMethod: func(c echo.Context) (map[string]interface{}, error) {
					return map[string]interface{}{
						"path":   c.Request().URL.Path,
						"method": c.Request().Method,
					}, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &EchoMiddleware{
				Config: tt.fields.Config,
				InputCreationMethod:   tt.fields.InputCreationMethod,
			}
			h.Use()
		})
	}
}

func TestNewEchoMiddleware(t *testing.T) {
	type args struct {
		cfg  *config.Config
		inputCreationMethod EchoInputCreationMethod
	}
	tests := []struct {
		name    string
		args    args
		want    *EchoMiddleware
		wantErr bool
	}{
		{
			name: "Test NewEchoMiddleware",
			args: args{
				cfg: &config.Config{
					Policy: "policy",
					Query:  "data.query",
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				inputCreationMethod: func(c echo.Context) (map[string]interface{}, error) {
					return map[string]interface{}{
						"path":   c.Request().URL.Path,
						"method": c.Request().Method,
					}, nil
				},
			},
			want: &EchoMiddleware{
				Config: &config.Config{
					Policy: "policy",
					Query:  "data.query",
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				InputCreationMethod: func(c echo.Context) (map[string]interface{}, error) {
					return map[string]interface{}{
						"path":   c.Request().URL.Path,
						"method": c.Request().Method,
					}, nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEchoMiddleware(tt.args.cfg, tt.args.inputCreationMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEchoMiddleware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
